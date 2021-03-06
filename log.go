package suexec

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

type Log struct {
	log *os.File
}

func NewLog(path string) *Log {
	log, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "suexec failure: could not open log file\n")
		fmt.Fprintf(os.Stderr, "fopen %s\n", err)
		os.Exit(1)
	}

	return &Log{log: log}
}

func (self *Log) errOutput(is_error bool, format string, args ...interface{}) {
	if is_error {
		fmt.Fprintf(os.Stderr, "suexec policy violation: see suexec log for more details\n")
	}

	t := time.Now()
	fmt.Fprintf(self.log, "[%d-%.2d-%.2d %.2d:%.2d:%.2d]: ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Fprintf(self.log, format, args...)
}

func (self *Log) LogErr(format string, args ...interface{}) {
	self.errOutput(true, format, args...)
}

func (self *Log) LogNoErr(format string, args ...interface{}) {
	self.errOutput(false, format, args...)
}

func (self *Log) SetCloseOnExec() (is_success bool) {

	is_success = true
	if _, _, errno := syscall.Syscall(syscall.SYS_FCNTL, self.log.Fd(), syscall.F_SETFD, syscall.FD_CLOEXEC); errno != 0 {
		self.LogErr("error: can't set close-on-exec flag")
		is_success = false
	}

	return
}
