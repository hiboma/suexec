package group

import (
	"os/user"
	"regexp"
)

func LookupGidAndName(groupname string) (gid string, actual_gname string, err error) {

	r, err := regexp.Compile(`^\d+$`)
	if r.MatchString(groupname) {
		gid = groupname
		actual_gname = groupname
		return
	}

	gr, err := user.LookupGroup(groupname)
	if err != nil {
		return
	}
	gid = gr.Gid
	actual_gname = gr.Name

	return
}
