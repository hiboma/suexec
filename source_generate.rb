#!/usr/bin/env ruby
# -*- encoding: utf-8 -*-

AP_HTTPD_USER = ENV["AP_HTTPD_USER"] || "vagrant"
AP_UID_MIN    = ENV["AP_UID_MIN"]    || 501
AP_GID_MIN    = ENV["AP_GID_MIN"]    || 501
AP_DOC_ROOT   = ENV["AP_DOC_ROOT"]   || "/vagrant/misc"
AP_LOG_EXEC   = ENV["AP_LOG_EXEC"]   || "/tmp/suexec.log"

puts <<"...."
package suexec

import (
	"fmt"
)

const AP_HTTPD_USER = "#{AP_HTTPD_USER}"
const AP_UID_MIN = #{AP_UID_MIN}
const AP_GID_MIN = #{AP_GID_MIN}
const AP_DOC_ROOT = "#{AP_DOC_ROOT}"
const AP_LOG_EXEC = "#{AP_LOG_EXEC}"

func PrintConstants() {
	fmt.Printf(" -D AP_DOC_ROOT=%s\\n", AP_DOC_ROOT)
	fmt.Printf(" -D AP_GID_MIN=%d\\n", AP_GID_MIN)
	fmt.Printf(" -D AP_HTTPD_USER=%s\\n", AP_HTTPD_USER)
	fmt.Printf(" -D AP_LOG_EXEC=%s\\n", AP_LOG_EXEC)
	fmt.Printf(" -D AP_SAFE_PATH=%s\\n", AP_SAFE_PATH)
	fmt.Printf(" -D AP_UID_MIN=%d\\n", AP_UID_MIN)
}
....
