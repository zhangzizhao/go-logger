package new

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

const (
	KUnit = 256
)

type Severity int

var (
	rotateNum int
	maxSize   uint64
	fall      bool
	console   bool
	logLevel  Severity
)

const (
	FATAL Severity = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

var severityName = []string{
	FATAL:   "FATAL",
	ERROR:   "ERROR",
	WARNING: "WARNING",
	INFO:    "INFO",
	DEBUG:   "DEBUG",
}

const (
	numSeverity = 5 // param ?
)

type LogDir struct {
	Dir string

	FlushInterval int // s
	BufferSize    int // k
	_suffix       int
	isCover       bool
	_date         *time.Time

	Backend *FileBackend
}

type FileBackend struct {
	mu    sync.Mutex
	dir   string //directory for log files
	files [numSeverity]syncBuffer
}

type syncBuffer struct {
	*bufio.Writer
	file     *os.File
	count    uint64
	cur      int
	filePath string
}

func SetFall(v bool) {
	fall = v
}

func SetConsole(v bool) {
	console = v
}

func SetLogLevel(v interface{}) {
	if s, ok := v.(Severity); ok {
		logLevel = s
	} else {
		if s, ok := v.(string); ok {
			for i, name := range severityName {
				if name == s {
					logLevel = Severity(i)
				}
			}
		}
	}
}

func (f *LogDir) pathInit() error {
	if _, err := os.Stat(f.Dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(f.Dir, 0755)
		} else {
			return err
		}
	}
	return nil
}

func (f *LogDir) Init() error {
	if f.FlushInterval == 0 {
		f.FlushInterval = 3
	}
	if f.BufferSize == 0 {
		f.BufferSize = 256
	}
	if err := os.MkdirAll(f.Dir, 0755); err != nil {
		return err
	}
	var fb FileBackend
	for i := 0; i < numSeverity; i++ {
		fileName := path.Join(f.Dir, severityName[i]+".log")
		logb, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		fb.files[i] = syncBuffer{Writer: bufio.NewWriterSize(logb, f.BufferSize*KUnit), file: logb, filePath: fileName}
	}
	// go fb.monitorFiles()
	f.Backend = &fb
	go f.Backend.flushDaemon(f.FlushInterval)
	return nil
}

func (self *FileBackend) Log(s Severity, msg []byte) {
	self.mu.Lock()
	switch s {
	case FATAL:
		self.files[FATAL].write(msg)
	case ERROR:
		self.files[ERROR].write(msg)
	case WARNING:
		self.files[WARNING].write(msg)
	case INFO:
		self.files[INFO].write(msg)
	case DEBUG:
		self.files[DEBUG].write(msg)
	}
	if fall && s < INFO {
		self.files[INFO].write(msg)
	}
	self.mu.Unlock()
	if s == FATAL {
		self.Flush()
	}
	//////  move to log.go?
	if console {
		go os.Stdout.Write(msg)
		// os.Stdout.Write([]byte("\n"))
	}
	//////

}
func (self *FileBackend) flushDaemon(interval int) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		self.Flush()
	}
}

func (self *FileBackend) close() {
	self.Flush()
}

//  完成切割，换fb（lock住flush）
// func (self *FileBackend) monitorFiles() {
// }

func (self *syncBuffer) Sync() error {
	return self.file.Sync()
}

func (self *FileBackend) Flush() {
	self.mu.Lock()
	defer self.mu.Unlock()
	for i := 0; i < numSeverity; i++ {
		self.files[i].Flush()
		self.files[i].Sync()
	}
}

func (self *syncBuffer) write(b []byte) { // write test when cur or rotate
	if maxSize > 0 && rotateNum > 0 && self.count+uint64(len(b)) >= maxSize {
		os.Rename(self.filePath, self.filePath+fmt.Sprintf(".%03d", self.cur))
		self.cur++
		if self.cur >= rotateNum {
			self.cur = 0
		}
		self.count = 0
	}
	self.count += uint64(len(b))
	self.Writer.Write(b)
}
