package logx

import (
	"compress/gzip"
	"fmt"
	"github.com/bhmy-shm/gofks/core/errorx"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	hoursPerDay     = 24
	bufferSize      = 100
	defaultDirMode  = 0o755
	defaultFileMode = 0o600
)

type RotateRule interface {
	BackupFileName() string
	MarkRotated()
	OutDatedFiles() []string
	ShallRotate() bool
}

// RotateLogger 结构体是一个日志切割的实现。
// 需要实现 io.{Closer,Writer} 接口方法，从而实现写和关操作。
type RotateLogger struct {
	filename  string
	backup    string
	fp        *os.File   //文件指针
	rule      RotateRule //日志切割规则 rule
	channel   chan []byte
	done      chan struct{}
	compress  bool
	waitGroup sync.WaitGroup
	closeOnce sync.Once
}

// DailyRotateRule 日志规则实现(按天)
type DailyRotateRule struct {
	gzip        bool   //是否对轮转的日志文件进行gzip压缩
	days        int    //保存日志文件的天数
	rotatedTime string //上次日志轮转的时间
	filename    string //日志文件名
	delimiter   string //日志文件名和时间部分之间的分隔符
}

func DefaultRotateRule(filename, delimiter string, days int, gzip bool) RotateRule {
	return &DailyRotateRule{
		rotatedTime: getNowDate(),
		filename:    filename,
		delimiter:   delimiter,
		days:        days,
		gzip:        gzip,
	}
}

// BackupFileName 形成备份文件的文件名
// 使用 filename、delimiter 和当前日期（通过调用 getNowDate() 函数获取）按照一定的格式组合起来。
func (r *DailyRotateRule) BackupFileName() string {
	return fmt.Sprintf("%s%s%s", r.filename, r.delimiter, getNowDate())
}

// MarkRotated 调用 getNowDate() 函数来获取当前的日期，并将其赋值给 r 的 rotatedTime 字段。
// 这个操作表示日志文件已经被轮转，并将轮转的时间更新为当前时间。
func (r *DailyRotateRule) MarkRotated() {
	r.rotatedTime = getNowDate()
}

// OutDatedFiles 用于返回超过保留天数的日志文件列表。
func (r *DailyRotateRule) OutDatedFiles() []string {

	//如果保留天数<=0，则没有列表直接返回nil
	if r.days <= 0 {
		return nil
	}

	//构建文件匹配的名称格式, 根据是否gzip压缩
	var pattern string
	if r.gzip {
		pattern = fmt.Sprintf("%s%s*.gz", r.filename, r.delimiter)
	} else {
		pattern = fmt.Sprintf("%s%s*", r.filename, r.delimiter)
	}

	//使用 filepath.Glob() 函数根据模式查找匹配的文件，并将结果存储在 files 切片中。
	files, err := filepath.Glob(pattern)
	if err != nil {
		Errorf("failed to delete outdated log files, error: %s", err)
		return nil
	}

	//计算出保留天数的边界日期，并将其与 filename 和 delimiter 拼接起来，得到一个边界文件名
	var buf strings.Builder
	boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay*r.days)).Format(dateFormat)
	fmt.Fprintf(&buf, "%s%s%s", r.filename, r.delimiter, boundary)
	if r.gzip {
		buf.WriteString(".gz")
	}
	boundaryFile := buf.String()

	//遍历 files 切片，将早于边界文件名的日志文件添加到 outDates 切片中，并返回此切片
	var outDates []string
	for _, file := range files {
		if file < boundaryFile {
			outDates = append(outDates, file)
		}
	}

	return outDates
}

// ShallRotate 判断是否需要进行切割，通过日期时间判断
// @response true 需要切割
// @response false 不需要切割
func (r *DailyRotateRule) ShallRotate() bool {
	return len(r.rotatedTime) > 0 && getNowDate() != r.rotatedTime
}

// ================= RotateLogger 真正写入日志文件的Rotate实现 ==================

func NewLogger(filename string, rule RotateRule, compress bool) (*RotateLogger, error) {
	l := &RotateLogger{
		filename:  filename,
		channel:   make(chan []byte, bufferSize),
		done:      make(chan struct{}),
		rule:      rule,
		compress:  compress,
		waitGroup: sync.WaitGroup{},
	}
	if err := l.buildFile(); err != nil {
		return nil, err
	}

	//注意：每new一个会开启协程waitGroup，但不会等待 wait结束。
	l.startWorker()
	return l, nil
}

// Close 关闭文件 RotateLogger，trait io.Closer
func (l *RotateLogger) Close() error {
	var err error
	l.closeOnce.Do(func() {
		close(l.done)
		l.waitGroup.Wait()

		// 刷盘操作，关闭前将文件的当前内容同步到稳定存储空间中
		if err = l.fp.Sync(); err != nil {
			return
		}
		err = l.fp.Close()
	})
	return err
}

// Writer 写入文件 RotateLogger, trait io.Writer
func (l *RotateLogger) Write(data []byte) (int, error) {
	select {
	case l.channel <- data:
		fmt.Println("Write data:", string(data))
		return len(data), nil
	case <-l.done:
		log.Println(string(data))
		return 0, errorx.ErrCodeLogFileClosed
	}
}

// buildFile 根据Logger 配置参数，在指定目录下生成文件
func (l *RotateLogger) buildFile() error {

	// 生成备份文件名称
	l.backup = l.rule.BackupFileName()

	//判断是否存在目录文件
	if _, err := os.Stat(l.filename); err != nil {

		// 生成文件目录
		basePath := path.Dir(l.filename) //仅返回目录路径，去掉具体文件名
		if _, err = os.Stat(basePath); err != nil {
			if err = os.MkdirAll(basePath, defaultDirMode); err != nil {
				return err
			}
		}

		// 创建具体文件
		if l.fp, err = os.Create(l.filename); err != nil {
			return err
		}
	} else if l.fp, err = os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, defaultFileMode); err != nil {
		return err
	}

	return nil
}

// startWorker 开启协程，通过channel 执行写入日志文件任务
func (l *RotateLogger) startWorker() {
	l.waitGroup.Add(1)

	go func() {
		defer l.waitGroup.Done()

		for {
			select {
			case event := <-l.channel:
				l.write(event)
			case <-l.done:
				return
			}
		}
	}()
}

// write 向文件写入内容
func (l *RotateLogger) write(v []byte) {

	//写入时判断是否需要切割
	if l.rule.ShallRotate() {
		if err := l.rotate(); err != nil {
			log.Println(err)
		} else {
			l.rule.MarkRotated() //获取新的日志文件日期时间
		}
	}

	// 正常写入文件
	if l.fp != nil {
		l.fp.Write(v)
	}
}

// getBackupFilename 获取当前logger 的备份文件名称
func (l *RotateLogger) getBackupFilename() string {
	if len(l.backup) == 0 {
		return l.rule.BackupFileName()
	}

	return l.backup
}

// rotate 用于执行日志切割操作。
func (l *RotateLogger) rotate() error {

	//关闭当前的文件指针 fp，然后判断是否需要备份当前日志文件。
	if l.fp != nil {
		err := l.fp.Close()
		l.fp = nil
		if err != nil {
			return err
		}
	}

	//
	_, err := os.Stat(l.filename)
	if err == nil && len(l.backup) > 0 {
		//如果需要备份，则将当前日志文件重命名为备份文件，并调用 postRotate() 方法来处理备份文件。
		backupFilename := l.getBackupFilename()
		err = os.Rename(l.filename, backupFilename)
		if err != nil {
			return err
		}

		l.postRotate(backupFilename)
	}

	//最后，重新创建一个新的日志文件。
	l.backup = l.rule.BackupFileName()
	if l.fp, err = os.Create(l.filename); err != nil {
		return err
	}

	return nil
}

// ================ gzip 压缩操作 ===================

// postRotate 异步方式处理日志文件压缩
func (l *RotateLogger) postRotate(file string) {
	go func() {
		// we cannot use threading.GoSafe here, because of import cycle.
		l.maybeCompressFile(file)
		l.maybeDeleteOutdatedFiles()
	}()
}

// maybeCompressFile 方法用于压缩备份文件
func (l *RotateLogger) maybeCompressFile(file string) {
	if !l.compress {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			Error(r)
		}
	}()

	if _, err := os.Stat(file); err != nil {
		// file not exists or other error, ignore compression
		return
	}

	compressLogFile(file)
}

// maybeDeleteOutdatedFiles 方法用于删除过时的备份文件
func (l *RotateLogger) maybeDeleteOutdatedFiles() {

	// 获取备份文件列表，然后依次进行删除
	files := l.rule.OutDatedFiles()
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			Errorf("failed to remove outdated file: %s", file)
		}
	}
}

// compressLogFile 方法用于压缩日志文件
func compressLogFile(file string) {
	start := time.Now()
	Infof("compressing log file: %s", file)
	if err := gzipFile(file); err != nil {
		Errorf("compress error: %s", err)
	} else {
		Infof("compressed log file: %s, took %s", file, time.Since(start))
	}
}

// gzipFile() 方法用于对文件进行压缩
func gzipFile(file string) error {

	//首先打开源文件和目标文件
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(fmt.Sprintf("%s.gz", file))
	if err != nil {
		return err
	}
	defer out.Close()

	//使用 gzip.NewWriter() 创建一个压缩写入器
	w := gzip.NewWriter(out)

	//将源文件内容复制到压缩写入器中，最后关闭压缩写入器和源文件
	if _, err = io.Copy(w, in); err != nil {
		return err
	} else if err = w.Close(); err != nil {
		return err
	}

	//删除源文件
	return os.Remove(file)
}
