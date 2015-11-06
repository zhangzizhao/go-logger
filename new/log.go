package new

import (
	"fmt"
)

// LogDir.Backend.Log(log.ERROR, []byte("test"))

func (l *LogDir) Close() {
	l.Backend.close()
}

// func (l *LogDir) LogToStderr() {
// 	l.logToStderr = true
// }

func getMsg(args ...interface{}) []byte {
	//  add header and  check end
	return []byte(fmt.Sprint(args...))
}

func (l *LogDir) Debug(args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(DEBUG, msg)
}

func (l *LogDir) Debugf(format string, args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(DEBUG, msg)
}

func (l *LogDir) Info(args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(INFO, msg)
}

func (l *LogDir) Infof(format string, args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(INFO, msg)
}

func (l *LogDir) Warn(args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(WARNING, msg)
}

func (l *LogDir) Warnf(format string, args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(WARNING, msg)
}

func (l *LogDir) Error(args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(ERROR, msg)
}

func (l *LogDir) Errorf(format string, args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(ERROR, msg)
}

func (l *LogDir) Fatal(args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(FATAL, msg)
}

func (l *LogDir) Fatalf(format string, args ...interface{}) {
	msg := getMsg(args...)
	l.Backend.Log(FATAL, msg)
}
