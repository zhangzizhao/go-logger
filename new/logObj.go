package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

const (
	bufferSize = 256 * 1024 // param ?
)

type Severity int

var (
	flushInterval time.Duration = time.Second * 3 // param ?
	rotateNum     int
	maxSize       uint64
	fall          bool
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

type stdBackend struct{}

func (self *stdBackend) Log(s Severity, msg []byte) {
	os.Stdout.Write(msg)
}

type _LogDir struct {
	Dir string

	_suffix int
	isCover bool
	_date   *time.Time

	Backend *FileBackend
	// mu       *sync.RWMutex
	// f        *os.File
	// writer   *bufio.Writer
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

func (f *_LogDir) pathInit() error {
	if _, err := os.Stat(f.Dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(f.Dir, 0755)
		} else {
			return err
		}
	}
	return nil
}

func (f *_LogDir) Init() error {
	if err := os.MkdirAll(f.Dir, 0755); err != nil {
		return err
	}
	var fb FileBackend
	for i := 0; i < numSeverity; i++ {
		fileName := path.Join(f.Dir, severityName[i]+".log")
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		fb.files[i] = syncBuffer{Writer: bufio.NewWriterSize(f, bufferSize), file: f, filePath: fileName}
	}
	// go fb.monitorFiles()
	f.Backend = &fb
	go f.Backend.flushDaemon()
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
	// if fall && s < INFO {
	// 	self.files[INFO].write(msg)
	// }
	self.mu.Unlock()
	if s == FATAL {
		self.Flush()
	}
}
func (self *FileBackend) flushDaemon() {
	for {
		time.Sleep(flushInterval)
		self.Flush()
	}
}

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

func main() {
	log := _LogDir{
		Dir: "./log",
	}
	log.Init()
	log.Backend.Log(ERROR, []byte("test"))
	select {}
}
