package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	logrus "github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Log flag
var (
	// Log default
	Log = New() // default
	// RotationCount count
	RotationCount uint = 365 //
)

// Logger log struct
type Logger struct {
	*logrus.Logger
	// local
	LogFlag      int // 日志标签 非DEBUG方法
	LogFlagDebug int // 日志 DEBUG方法
}

// SetFlags log normal
func SetFlags(flg int) {
	Log.SetFlags(flg)
}

// SetLevel log normal
func SetLevel(level uint32) {
	Log.SetLevel(level)
}

// Debug debug
func Debug(args ...interface{}) {
	if Log.Level >= logrus.DebugLevel {
		Log.EntryWith(Log.LogFlag).Debug(args...)
	}
}

// Debugf debug
func Debugf(format string, args ...interface{}) {
	if Log.Level >= logrus.DebugLevel {
		Log.EntryWith(Log.LogFlag).Debugf(format, args...)
	}
}

// Debugln debug
func Debugln(args ...interface{}) {
	if Log.Level >= logrus.DebugLevel {
		Log.EntryWith(Log.LogFlag).Debugln(args...)
	}
}

// Info info
func Info(args ...interface{}) {
	if Log.Level >= logrus.InfoLevel {
		Log.EntryWith(Log.LogFlag).Info(args...)
	}
}

// Infof info
func Infof(format string, args ...interface{}) {
	if Log.Level >= logrus.InfoLevel {
		Log.EntryWith(Log.LogFlag).Infof(format, args...)
	}
}

// Infoln info
func Infoln(args ...interface{}) {
	if Log.Level >= logrus.InfoLevel {
		Log.EntryWith(Log.LogFlag).Infoln(args...)
	}
}

// Warn warn
func Warn(args ...interface{}) {
	if Log.Level >= logrus.WarnLevel {
		Log.EntryWith(Log.LogFlag).Warn(args...)
	}
}

// Warnf warn
func Warnf(format string, args ...interface{}) {
	if Log.Level >= logrus.WarnLevel {
		Log.EntryWith(Log.LogFlag).Warnf(format, args...)
	}
}

// Warnln warn
func Warnln(args ...interface{}) {
	if Log.Level >= logrus.WarnLevel {
		Log.EntryWith(Log.LogFlag).Warnln(args...)
	}
}

// Error error
func Error(args ...interface{}) {
	if Log.Level >= logrus.ErrorLevel {
		Log.EntryWith(Log.LogFlag).Error(args...)
	}
}

// Errorf error
func Errorf(format string, args ...interface{}) {
	if Log.Level >= logrus.ErrorLevel {
		Log.EntryWith(Log.LogFlag).Errorf(format, args...)
	}
}

// Errorln error
func Errorln(args ...interface{}) {
	if Log.Level >= logrus.ErrorLevel {
		Log.EntryWith(Log.LogFlag).Errorln(args...)
	}
}

// Fatal fatal
func Fatal(args ...interface{}) {
	if Log.Level >= logrus.FatalLevel {
		Log.EntryWith(Log.LogFlag).Fatal(args...)
	}
}

// Fatalf fatal
func Fatalf(format string, args ...interface{}) {
	if Log.Level >= logrus.FatalLevel {
		Log.EntryWith(Log.LogFlag).Fatalf(format, args...)
	}
}

// Fatalln fatal
func Fatalln(args ...interface{}) {
	if Log.Level >= logrus.FatalLevel {
		Log.EntryWith(Log.LogFlag).Fatalln(args...)
	}
}

// Panic panic
func Panic(args ...interface{}) {
	Log.EntryWith(Log.LogFlag).Panic(args...)
}

// Panicf panic
func Panicf(format string, args ...interface{}) {
	Log.EntryWith(Log.LogFlag).Panicf(format, args...)
}

// Panicln panic
func Panicln(args ...interface{}) {
	Log.EntryWith(Log.LogFlag).Panicln(args...)
}

// Print print
func Print(args ...interface{}) {
	logrus.NewEntry(Log.Logger).Print(args...)
}

// Printf print
func Printf(format string, args ...interface{}) {
	logrus.NewEntry(Log.Logger).Printf(format, args...)
}

// Println print
func Println(args ...interface{}) {
	logrus.NewEntry(Log.Logger).Println(args...)
}

// SetFlags log normal
func (l *Logger) SetFlags(flg int) {
	l.LogFlag |= flg
	l.LogFlagDebug = l.LogFlag
}

// SetLevel log normal
func (l *Logger) SetLevel(level uint32) {
	l.Level = logrus.Level(level)
}

// GetLevel log normal
func (l *Logger) GetLevel() uint32 {
	return uint32(l.Level)
}

// Upon up on
func (l *Logger) Upon(level uint32) (ret bool) {
	if uint32(l.Level) >= level {
		ret = true
	}
	return
}

// Debug debug
func (l *Logger) Debug(args ...interface{}) {
	if l.Level >= logrus.DebugLevel {
		l.EntryWith(l.LogFlagDebug).Debug(args...)
	}
}

// Debugf debug
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Level >= logrus.DebugLevel {
		l.EntryWith(l.LogFlagDebug).Debugf(format, args...)
	}
}

// Debugln debug
func (l *Logger) Debugln(args ...interface{}) {
	if l.Level >= logrus.DebugLevel {
		l.EntryWith(l.LogFlagDebug).Debugln(args...)
	}
}

// Info info
func (l *Logger) Info(args ...interface{}) {
	if l.Level >= logrus.InfoLevel {
		l.EntryWith(l.LogFlag).Info(args...)
	}
}

// Infof info
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.Level >= logrus.InfoLevel {
		l.EntryWith(l.LogFlag).Infof(format, args...)
	}
}

// Infoln info
func (l *Logger) Infoln(args ...interface{}) {
	if l.Level >= logrus.InfoLevel {
		l.EntryWith(l.LogFlag).Infoln(args...)
	}
}

// Warn warn
func (l *Logger) Warn(args ...interface{}) {
	if l.Level >= logrus.WarnLevel {
		l.EntryWith(l.LogFlag).Warn(args...)
	}
}

// Warnf warn
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.Level >= logrus.WarnLevel {
		l.EntryWith(l.LogFlag).Warnf(format, args...)
	}
}

// Warnln warn
func (l *Logger) Warnln(args ...interface{}) {
	if l.Level >= logrus.WarnLevel {
		l.EntryWith(l.LogFlag).Warnln(args...)
	}
}

// Error error
func (l *Logger) Error(args ...interface{}) {
	if l.Level >= logrus.ErrorLevel {
		l.EntryWith(l.LogFlag).Error((args)...)
	}
}

// Errorf error
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.Level >= logrus.ErrorLevel {
		l.EntryWith(l.LogFlag).Errorf(format, (args)...)
	}
}

// Errorln error
func (l *Logger) Errorln(args ...interface{}) {
	if l.Level >= logrus.ErrorLevel {
		l.EntryWith(l.LogFlag).Errorln((args)...)
	}
}

// Print print
func (l *Logger) Print(args ...interface{}) {
	l.EntryWith(l.LogFlag).Print(args...)
}

// Printf print
func (l *Logger) Printf(format string, args ...interface{}) {
	l.EntryWith(l.LogFlag).Printf(format, args...)
}

// Println print
func (l *Logger) Println(args ...interface{}) {
	l.EntryWith(l.LogFlag).Println(args...)
}

// Fatal fatal
func (l *Logger) Fatal(args ...interface{}) {
	if l.Level >= logrus.FatalLevel {
		l.EntryWith(l.LogFlag).Fatal(args...)
	}
}

// Fatalf fatal
func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.Level >= logrus.FatalLevel {
		l.EntryWith(l.LogFlag).Fatalf(format, args...)
	}
}

// Fatalln fatal
func (l *Logger) Fatalln(args ...interface{}) {
	if l.Level >= logrus.FatalLevel {
		l.EntryWith(l.LogFlag).Fatalln(args...)
	}
}

// Panic panic
func (l *Logger) Panic(args ...interface{}) {
	if l.Level >= logrus.PanicLevel {
		l.EntryWith(l.LogFlag).Panic(args...)
	}
}

// Panicf panic
func (l *Logger) Panicf(format string, args ...interface{}) {
	if l.Level >= logrus.PanicLevel {
		l.EntryWith(l.LogFlag).Panicf(format, args...)
	}
}

// Panicln panic
func (l *Logger) Panicln(args ...interface{}) {
	if l.Level >= logrus.PanicLevel {
		l.EntryWith(l.LogFlag).Panicln(args...)
	}
}

// EntryWith 格式化输出
func (l *Logger) EntryWith(flg int) *logrus.Entry {
	if flg&(log.Lshortfile|log.Llongfile) != 0 {
		if pc, file, line, ok := runtime.Caller(2); ok {
			var (
				_fnName    = runtime.FuncForPC(pc).Name()
				_fnNameLis = strings.Split(_fnName, ".")
			)
			fnName := _fnNameLis[len(_fnNameLis)-1]
			return l.Logger.WithFields(logrus.Fields{
				"file": file,
				"line": line,
				"func": fnName,
			})
		}
	}

	return logrus.NewEntry(l.Logger)
}

// Recover panic 错误捕捉
func (l *Logger) Recover() {
	r := recover()
	if r == nil {
		return
	}
	var (
		pc     [16]uintptr
		n      = runtime.Callers(3, pc[:])
		fields logrus.Fields
	)
	for _, _pc := range pc[:n] {
		fn := runtime.FuncForPC(_pc)
		if fn == nil {
			continue
		}
		_fnName := fn.Name()
		if strings.HasPrefix(_fnName, "runtime.") {
			continue
		}
		file, line := fn.FileLine(_pc)

		//
		var (
			_fnNameDir = strings.Split(_fnName, "/")
			_fnNameLis = strings.Split(_fnName, ".")
			_fnNameSrc string
		)
		if len(_fnNameDir) > 1 {
			_fnNameSrc = _fnNameDir[0] + "/" + _fnNameDir[1] + "/"
		} else {
			_fnNameSrc = _fnNameDir[0]
		}
		fnName := _fnNameLis[len(_fnNameLis)-1]

		// file
		_pcLis := strings.Split(file, _fnNameSrc)
		filePath := strings.Join(_pcLis, "")
		fields = logrus.Fields{
			"file": filePath,
			"line": line,
			"func": fnName,
		}
		goto PRINT
	}
PRINT:
	l.Logger.WithFields(fields).Errorln(r)
}

// Close 关闭
func (l *Logger) Close() error {
	if l.Out != nil {
		if w, ok := l.Out.(io.WriteCloser); ok {
			return w.Close()
		}
	}
	return nil
}

// Copy 复制
func (l *Logger) Copy() (r *Logger) {
	r = New()
	r.SetLevel(uint32(l.Level))
	r.Out = l.Out
	return
}

func New() *Logger {
	d := &Logger{}
	d.Logger = logrus.New()
	// default
	d.SetFlags(log.Llongfile)
	d.SetLevel(uint32(logrus.DebugLevel))
	d.Out = os.Stderr
	return d
}

// NewLogFile new log file
func NewLogFile(logPath string) (d *Logger) {
	var (
		//f   *os.File
		rf, rfErr *rotatelogs.RotateLogs
		err       error
	)
	d = New()

	// ensure director
	_dir := filepath.Dir(logPath)
	if _, _err := os.Stat(_dir); os.IsNotExist(_err) {
		if err = os.MkdirAll(_dir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	// log file(s)
	if rf, err = rotatelogs.New(
		logPath+".%Y-%m-%d.log",
		//rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		// rotatelogs.WithRotationCount(RotationCount),
	); err == nil {
		rfErr, _ = rotatelogs.New(
			logPath+".%Y-%m-%d.error",
			//rotatelogs.WithLinkName(logPath),
			rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
			rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
			// rotatelogs.WithRotationCount(RotationCount),
		)
		d.Hooks.Add(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.TraceLevel: rf,
				logrus.DebugLevel: rf,
				logrus.InfoLevel:  rf,
				logrus.WarnLevel:  rf,
				logrus.ErrorLevel: rfErr,
				logrus.FatalLevel: rf,
				logrus.PanicLevel: rf,
			},
			&logrus.JSONFormatter{},
		))

	} else {
		logrus.Warnln(err)
	}

	return
}

// HookSelf hook self 钩子
type HookSelf struct {
	writers   lfshook.WriterMap
	levels    []logrus.Level
	lock      *sync.Mutex
	formatter logrus.Formatter
}

// Levels level
func (h *HookSelf) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *HookSelf) SetFormat(logger *logrus.Logger, format logrus.Formatter) *HookSelf {
	logger.SetFormatter(format)
	return h
}

// Fire fire
func (h *HookSelf) Fire(entry *logrus.Entry) error {
	var (
		writer io.Writer
		msg    []byte
		err    error
	)

	// use our formatter instead of entry.String()
	msg, err = h.formatter.Format(entry)

	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	_, err = writer.Write([]byte(msg))
	return nil
}

// PanicRecover 统一处理panic
func PanicRecover(logger *Logger) {
	r := recover()
	if r == nil {
		return
	}
	if logger == nil {
		logger = New()
	}
	logger.Errorf(`[panic-recover] "%s" %v`, panicIdentify(), r)
}

// PanicRecoverError 统一处理panic, 并更新error
func PanicRecoverError(logger *Logger, err *error) {
	r := recover()
	if r == nil {
		return
	}
	if logger == nil {
		logger = New()
	}
	logger.Errorf(`[panic-recover] "%s" %v`, panicIdentify(), r)
	return
}

// panicIdentify  定位panic位置 参考自: https://gist.github.com/swdunlop/9629168
func panicIdentify() string {
	var (
		pc [16]uintptr
		n  = runtime.Callers(3, pc[:])
	)
	for _, _pc := range pc[:n] {
		fn := runtime.FuncForPC(_pc)
		if fn == nil {
			continue
		}
		_fnName := fn.Name()
		if strings.HasPrefix(_fnName, "runtime.") {
			continue
		}
		file, line := fn.FileLine(_pc)

		//
		var (
			_fnNameDir = strings.Split(_fnName, "/")
			_fnNameLis = strings.Split(_fnName, ".")
			_fnNameSrc string
		)
		if len(_fnNameDir) > 1 {
			_fnNameSrc = _fnNameDir[0] + "/" + _fnNameDir[1] + "/"
		} else {
			_fnNameSrc = _fnNameDir[0]
		}
		fnName := _fnNameLis[len(_fnNameLis)-1]

		// file
		_pcLis := strings.Split(file, _fnNameSrc)
		filePath := strings.Join(_pcLis, "")
		return fmt.Sprintf("%s:%d|%s", filePath, line, fnName)
	}
	return "unknown"
}
