package pkg

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

type (
	WatcherInter interface {
		Next(config *Config) (*File, error)
		Stop() error
	}

	watcher struct {
		f    *File
		fw   *fsnotify.Watcher
		exit chan bool
	}
)

func newWatcher(f *File) (WatcherInter, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fw.Add(f.opts.path)

	return &watcher{
		f:    f,
		fw:   fw,
		exit: make(chan bool),
	}, nil
}

func (w *watcher) Next(conf *Config) (*File, error) {
	// is it closed?
	select {
	case <-w.exit:
		return nil, errorx.ErrCodeWatcherFileStop
	default:
	}

	select {
	case event, _ := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			// check existence of File, and add watch again
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		}

		//重新读取加载文件，并返回文件
		return w.f.Reload(conf), nil

	case err := <-w.fw.Errors:
		return nil, err
	case <-w.exit:
		return nil, errorx.ErrCodeWatcherFileStop
	}
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}

// ReadWatcher 循环监听配置文件修改
func ReadWatcher(f *File, conf *Config) {
	w, _ := f.Watch()

	newf, err := w.Next(conf)
	if err != nil {
		log.Fatalln("watch faile =", err)
	}

	//新旧替换(包含了文件和配置实例对象)
	f = newf

	if f.confMap != nil {
		ReadWatcher(f, conf)
	} else {
		log.Fatalln("监听后的内容出现异常")
	}
}
