package suexec

import (
	"fmt"
	"strings"
)

type Error struct {
	status  int
	message string
}

func NewError(status int, format string, args ...interface{}) *Error {
	return &Error{status, fmt.Sprintf(format, args...)}
}

func (self *Error) Status() int {
	return self.status
}

func (self *Error) Message() string {
	return self.message
}

const AP_HTTPD_USER = "hiroya"
const AP_UID_MIN = "500"
const AP_GID_MIN = "10"
const AP_DOC_ROOT = "/private/var/tmp/"
const AP_LOG_EXEC = "/private/var/tmp/suexec.log"

func IsUserdirEnabled(username string) bool {
	return strings.HasPrefix(username, "~")
}
