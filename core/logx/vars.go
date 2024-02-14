package logx

const (
	dateFormat          = "2006-01-02"
	timeFormat          = "2006-01-02T15:04:05"
	callerDepth         = 5   //调用栈的深度
	backupFileDelimiter = "-" //
)

const (
	// InfoLevel logs everything
	InfoLevel uint32 = iota
	// ErrorLevel includes errors
	ErrorLevel
	// StatLevel includes stat
	StatLevel
	// SevereLevel only log severe messages
	SevereLevel
)

const (
	infoFilename  = "info.log"
	errorFilename = "error.log"
	statFile      = "stat.log"

	flags     = 0x0
	cronFlags = 3
)

const (
	consoleMode = "console"
	fileMode    = "file"
)

const (
	levelAlert = "alert"
	levelInfo  = "info"
	levelError = "error"
	levelFatal = "fatal"
	levelStat  = "stat"
	levelSlow  = "slow"
)

const (
	jsonEncodingType = iota
	plainEncodingType
	jsonEncoding     = "json"
	plainEncoding    = "plain"
	plainEncodingSep = '\t'
)

const (
	callerKey    = "caller"
	contentKey   = "content"
	durationKey  = "duration"
	levelKey     = "level"
	spanKey      = "span"
	timestampKey = "@timestamp"
	traceKey     = "trace"
)
