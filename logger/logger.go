package logger

import (
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	logger *log.Logger
	mu     sync.Mutex
}

var Log *Logger

func NewLogger(logDir string) *Logger {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}
	logFile, err := os.OpenFile(logDir+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法打开日志文件:", err)
	}
	return &Logger{
		logger: log.New(logFile, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.SetPrefix("[INFO] ")
	l.logger.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(v) > 0 {
		l.logger.Printf("%s:%d: "+format+"\n", append([]interface{}{file, line}, v...)...)
	} else {
		l.logger.Printf("%s:%d: "+format+"\n", file, line)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.SetPrefix("[ERROR] ")
	l.logger.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(v) > 0 {
		l.logger.Printf("%s:%d: "+format+"\n", append([]interface{}{file, line}, v...)...)
	} else {
		l.logger.Printf("%s:%d: "+format+"\n", file, line)
	}
}

func (l *Logger) Close() {
	l.logger.Writer().(*os.File).Close()
}

func InitLog(logDir string) *Logger {
	logger := NewLogger(logDir)
	Log = logger
	return logger
}
