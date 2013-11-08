package suexec

import (
	"fmt"
)

const AP_HTTPD_USER = "vagrant"
const AP_UID_MIN = 501
const AP_GID_MIN = 501
const AP_DOC_ROOT = "/vagrant/misc"
const AP_LOG_EXEC = "/tmp/suexec.log"

func PrintConstants() {
	fmt.Printf(" -D AP_DOC_ROOT=%s\n", AP_DOC_ROOT)
	fmt.Printf(" -D AP_GID_MIN=%d\n", AP_GID_MIN)
	fmt.Printf(" -D AP_HTTPD_USER=%s\n", AP_HTTPD_USER)
	fmt.Printf(" -D AP_LOG_EXEC=%s\n", AP_LOG_EXEC)
	fmt.Printf(" -D AP_SAFE_PATH=%s\n", AP_SAFE_PATH)
	fmt.Printf(" -D AP_UID_MIN=%d\n", AP_UID_MIN)
}
