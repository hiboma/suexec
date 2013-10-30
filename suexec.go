package suexec

import (
	"fmt"
	"strings"
)

const AP_HTTPD_USER = "vagrant"
const AP_UID_MIN = "501"
const AP_GID_MIN = "501"
const AP_DOC_ROOT = "/vagrant/misc"
const AP_LOG_EXEC = "/tmp/suexec.log"

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

func IsUserdirEnabled(username string) bool {
	return strings.HasPrefix(username, "~")
}
