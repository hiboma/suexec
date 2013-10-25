package suexec

import (
	"strings"
)

const AP_HTTPD_USER = "apache"
const AP_UID_MIN = "500"
const AP_GID_MIN = "10"
const AP_DOC_ROOT = "/var/www"
const AP_LOG_EXEC = "/var/log/httpd/suexec.log"

func IsUserdirEnabled(username string) bool {
	return strings.HasPrefix(username, "~")
}

func IsValidCommand(cmd string) bool {
	if strings.HasPrefix(cmd, "/") ||
		strings.HasPrefix(cmd, "../") ||
		strings.Index(cmd, "/../") > 0 {
		return false
	}
	return true
}
