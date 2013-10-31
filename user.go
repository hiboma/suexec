package suexec

import (
	"os/user"
	"regexp"
	"strings"
)

func IsUserdirEnabled(username string) bool {
	return strings.HasPrefix(username, "~")
}

func Lookup(username string) (pw *user.User, err error) {

	r, err := regexp.Compile(`^\d+$`)
	if !r.MatchString(username) {
		pw, err = user.Lookup(username)
	} else {
		pw, err = user.LookupId(username)
	}

	return
}
