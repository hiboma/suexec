package suexec

import (
	"fmt"
)

type SuexecError struct {
	status  int
	message string
}

func NewSuexecError(status int, format string, args ...interface{}) *SuexecError {
	return &SuexecError{status, fmt.Sprintf(format, args...)}
}

func (self *SuexecError) Status() int {
	return self.status
}

func (self *SuexecError) Message() string {
	return self.message
}
