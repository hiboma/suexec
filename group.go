package suexec

import (
	"os/user"
	"regexp"
)

func LookupGroup(groupname string) (*user.Group, error) {

	var gr *user.Group
	var err error

	r := regexp.MustCompile(`^\d+$`)
	if r.MatchString(groupname) {
		gr, err = user.LookupGroupId(groupname)
		if err != nil {
			return nil, err
		}
	} else {
		gr, err = user.LookupGroup(groupname)
		if err != nil {
			return nil, err
		}
	}

	return gr, err
}
