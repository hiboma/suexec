package suexec

import (
	"os/user"
	"regexp"
)

func Lookup(username string) (pw *user.User, err error) {

	r, err := regexp.Compile(`^\d+$`)
	if !r.MatchString(username) {
		pw, err = user.Lookup(username)
	} else {
		pw, err = user.LookupId(username)
	}

	return
}
