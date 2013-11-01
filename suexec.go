package suexec

import (
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

type Param struct {
	args []string
	log  *Log
	uid  int
	cwd  string
}

func Suexec(p Param) *SuexecError {

	var userdir bool = false

	args := p.args
	cwd := p.cwd
	log := p.log
	original_uid := p.uid

	/*
	 * Start with a "clean" environment
	 */
	environ := CleanEnv()

	/*
	 * Check existence/validity of the UID of the user
	 * running this program.  Error out if invalid.
	 */
	pw, err := user.LookupId(strconv.Itoa(original_uid))
	if err != nil {
		return NewSuexecError(102, "crit: invalid uid: (%d) %s\n", original_uid, err)
	}
	/*
	 * See if this is a 'how were you compiled' request, and
	 * comply if so.
	 */
	if len(args) > 1 && args[1] == "-V" && pw.Uid == "0" {
		PrintConstants()
		return NewSuexecError(0, "")
	}
	/*
	 * If there are a proper number of arguments, set
	 * all of them to variables.  Otherwise, error out.
	 */
	if len(args) < 4 {
		return NewSuexecError(101, "too few arguments")
	}

	target_uname := args[1]
	target_gname := args[2]
	cmd := args[3]
	/*
	 * Check to see if the user running this program
	 * is the user allowed to do so as defined in
	 * suexec.h.  If not the allowed user, error out.
	 */
	if AP_HTTPD_USER != pw.Username {
		return NewSuexecError(103, "user mismatch (%s instead os %s)", pw.Username, AP_HTTPD_USER)
	}

	/*
	 * Check to see if this is a ~userdir request.  If
	 * so, set the flag, and remove the '~' from the
	 * target username.
	 */
	userdir = IsUserdirEnabled(target_uname)

	/*
	 * Error out if the target username is invalid.
	 */
	pw, err = Lookup(target_uname)
	if err != nil {
		return NewSuexecError(121, "invalid target user: (%s)\n", target_uname)
	}

	/*
	 * Error out if the target group name is invalid.
	 */
	gr, err := LookupGroup(target_gname)
	if err != nil {
		return NewSuexecError(106, "invalid target group name: (%s)\n", target_gname)
	}
	gid, err := strconv.Atoi(gr.Gid)
	if err != nil {
		log.LogErr("failed to strconv.Atoi: (%v)\n", err)
	}
	actual_gname := gr.Name

	/*
	 * Save these for later since initgroups will hose the struct
	 */
	uid, err := strconv.Atoi(pw.Uid)
	if err != nil {
		return NewSuexecError(255, "failed to strconv.Atoi: (%v)\n", err)
	}
	actual_uname := pw.Username
	//	target_homedir := pw.HomeDir

	/*
	 * Log the transaction here to be sure we have an open log
	 * before we setuid().
	 */
	log.LogNoErr("uid: (%s/%s) gid: (%s/%s) cmd: %s\n",
		target_uname, actual_uname,
		target_gname, actual_gname,
		cmd)

	/*
	 * Error out if attempt is made to execute as root or as
	 * a UID less than AP_UID_MIN.  Tsk tsk.
	 */
	if uid == 0 || uid < AP_UID_MIN {
		return NewSuexecError(107, "cannot run as forbidden uid (%d/%s)\n", uid, cmd)
	}

	/*
	 * Error out if attempt is made to execute as root group
	 * or as a GID less than AP_GID_MIN.  Tsk tsk.
	 */
	if gid == 0 || (gid < AP_GID_MIN) {
		return NewSuexecError(108, "cannot run as forbidden gid (%d/%s)\n", gid, cmd)
	}

	/*
	 * Change UID/GID here so that the following tests work over NFS.
	 *
	 * Initialize the group access list for the target user,
	 * and setgid() to the target group. If unsuccessful, error out.
	 */
	if err := syscall.Setgid(gid); err != nil {
		return NewSuexecError(109, "failed to setgid (%d: %s)\n", gid, cmd)
	}

	/*
	 * setuid() to the target user.  Error out on fail.
	 */
	if err := syscall.Setuid(uid); err != nil {
		return NewSuexecError(110, "failed to setuid (%d: %s)\n", uid, cmd)
	}

	/*
	 * Get the current working directory, as well as the proper
	 * document root (dependant upon whether or not it is a
	 * ~userdir request).  Error out if we cannot get either one,
	 * or if the current working directory is not in the docroot.
	 * Use chdir()s and getcwd()s to avoid problems with symlinked
	 * directories.  Yuck.
	 */
	var dwd string

	if userdir {
		/* todo */
	} else {
		/* oops */
		if err = os.Chdir(AP_DOC_ROOT); err != nil {
			return NewSuexecError(113, "cannot get docroot information (%s)\n", AP_DOC_ROOT)
		}
		dwd, err = os.Getwd()
		if err != nil {
			return NewSuexecError(113, "cannot get docroot information (%s)\n", AP_DOC_ROOT)
		}

		if err = os.Chdir(cwd); err != nil {
			return NewSuexecError(113, "cannot get docroot information (%s)\n", AP_DOC_ROOT)
		}
	}

	if !strings.HasPrefix(cwd, dwd) {
		return NewSuexecError(114, "command not in docroot (%s/%s)\n", cwd, cmd)
	}

	script, err := NewScript(cmd, cwd)
	if err != nil {
		return NewSuexecError(1, "%v\n", err)
	}

	if suexec_err := script.VerifyToSuexec(uid, gid); err != nil {
		return suexec_err
	}

	log.SetCloseOnExec()

	if err := syscall.Exec(cmd, args, environ); err != nil {
		return NewSuexecError(1, "(%d) %s: failed(%s)\n", err, err, cmd)
	}

	/*
	 * (I can't help myself...sorry.)
	 *
	 * Uh oh.  Still here.  Where's the kaboom?  There was supposed to be an
	 * EARTH-shattering kaboom!
	 *
	 * Oh well, log the failure and error out.
	 */
	return NewSuexecError(255, "(%d)%s: exec failed (%s)\n", err, err, cmd)
}
