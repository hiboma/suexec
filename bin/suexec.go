package main

import (
	"fmt"
	"github.com/hiboma/suexec"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var log os.File

func logErr(format string, args ...interface{}) {
	t := time.Now()
	fmt.Fprintf(os.Stderr, "%s", t.Year)
	fmt.Fprintf(os.Stderr, format, args...)
}

func logNoErr(format string, args ...interface{}) {

	//if !log {
	log, err := os.OpenFile(suexec.AP_LOG_EXEC, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "suexec failure: could not open log file\n")
		fmt.Fprintf(os.Stderr, "fopen %s\n", err)
		os.Exit(1)
	}
	//}

	t := time.Now()
	fmt.Fprintf(log, "[%d-%.2d-%.2d %.2d:%.2d:%.2d]: ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Fprintf(log, format, args...)
}

func main() {
	var userdir bool = false

	/*
	 * Start with a "clean" environment
	 */
	environ := suexec.CleanEnv()

	/*
	 * Check existence/validity of the UID of the user
	 * running this program.  Error out if invalid.
	 */
	pw, err := user.Current()
	if err != nil {
		logErr("crit: invalid uid: (%d) %s\n", os.Getuid(), err)
		os.Exit(1)
	}
	/*
	 * See if this is a 'how were you compiled' request, and
	 * comply if so.
	 */
	if len(os.Args) > 1 && os.Args[1] == "-V" && pw.Uid == "0" {
		fmt.Printf(" -D AP_DOC_ROOT=%s\n")
		os.Exit(0)
	}
	/*
	 * If there are a proper number of arguments, set
	 * all of them to variables.  Otherwise, error out.
	 */
	if len(os.Args) < 4 {
		fmt.Println("too few arguments")
		os.Exit(101)
	}

	target_uname := os.Args[1]
	target_gname := os.Args[2]
	cmd := os.Args[3]
	/*
	 * Check to see if the user running this program
	 * is the user allowed to do so as defined in
	 * suexec.h.  If not the allowed user, error out.
	 */
	if suexec.AP_HTTPD_USER != pw.Username {
		logErr("user mismatch (%s instead os %s)", pw.Username, suexec.AP_HTTPD_USER)
		os.Exit(103)
	}

	/*
	 * Check to see if this is a ~userdir request.  If
	 * so, set the flag, and remove the '~' from the
	 * target username.
	 */
	userdir = suexec.IsUserdirEnabled(target_uname)

	/*
	 * Error out if the target username is invalid.
	 */
	pw, err = suexec.Lookup(target_uname)
	if err != nil {
		logErr("invalid target user: (%s)\n", target_uname)
		os.Exit(121)
	}

	/*
	 * Error out if the target group name is invalid.
	 */
	gid, actual_gname, err := suexec.LookupGidAndName(target_gname)
	if err != nil {
		logErr("invalid target group name: (%s)\n", target_gname)
		os.Exit(106)
	}

	/*
	 * Save these for later since initgroups will hose the struct
	 */
	uid := pw.Uid
	actual_uname := pw.Username
	//	target_homedir := pw.HomeDir

	/*
	 * Log the transaction here to be sure we have an open log
	 * before we setuid().
	 */
	logNoErr("uid: (%s/%s) gid: (%s/%s) cmd: %s\n",
		target_uname, actual_uname,
		target_gname, actual_gname,
		cmd)

	/*
	 * Error out if attempt is made to execute as root or as
	 * a UID less than AP_UID_MIN.  Tsk tsk.
	 */
	if uid == "0" || uid < suexec.AP_UID_MIN {
		logErr("cannot run as forbidden uid (%d/%s)\n", uid, cmd)
		os.Exit(107)
	}

	/*
	 * Error out if attempt is made to execute as root group
	 * or as a GID less than AP_GID_MIN.  Tsk tsk.
	 */
	if gid == "0" || (gid < suexec.AP_GID_MIN) {
		logErr("cannot run as forbidden gid (%lu/%s)\n", gid, cmd)
		os.Exit(108)
	}

	/*
	 * Change UID/GID here so that the following tests work over NFS.
	 *
	 * Initialize the group access list for the target user,
	 * and setgid() to the target group. If unsuccessful, error out.
	 */
	gid_int, err := strconv.Atoi(gid)
	if err := syscall.Setgid(gid_int); err != nil {
		logErr("failed to setgid (%lu: %s)\n", gid, cmd)
		os.Exit(109)
	}

	/*
	 * setuid() to the target user.  Error out on fail.
	 */
	uid_int, err := strconv.Atoi(uid)
	if err := syscall.Setuid(uid_int); err != nil {
		logErr("failed to setuid (%d: %s)\n", uid, cmd)
		os.Exit(110)
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
	cwd, err := os.Getwd()
	if err != nil {
		logErr("cannot get current working directory\n")
		os.Exit(111)
	}

	if userdir {
		/* todo */
	} else {
		/* oops */
		err = os.Chdir(suexec.AP_DOC_ROOT)
		if err != nil {
			logErr("cannot get docroot information (%s)\n", suexec.AP_DOC_ROOT)
			os.Exit(113)
		}
		dwd, err = os.Getwd()
		if err != nil {
			logErr("cannot get docroot information (%s)\n", suexec.AP_DOC_ROOT)
			os.Exit(113)
		}

		err = os.Chdir(cwd)
		if err != nil {
			logErr("cannot get docroot information (%s)\n", suexec.AP_DOC_ROOT)
			os.Exit(113)
		}
	}

	if !strings.HasPrefix(cwd, dwd) {
		logErr("command not in docroot (%s/%s)\n", cwd, cmd)
		os.Exit(114)
	}

	script, err := suexec.NewScript(cmd, cwd)
	if err != nil {
		logErr("%v\n", err)
		os.Exit(1)
	}

	if err := script.VerifyToSuexec(uid_int, gid_int); err != nil {
		logErr(err.Message())
		os.Exit(err.Status())
	}

	if err := syscall.Exec(cmd, os.Args, environ); err != nil {
		logErr("(%d) %s: failed(%s)\n", err, err, cmd)
		os.Exit(1)
	}

	/*
	 * (I can't help myself...sorry.)
	 *
	 * Uh oh.  Still here.  Where's the kaboom?  There was supposed to be an
	 * EARTH-shattering kaboom!
	 *
	 * Oh well, log the failure and error out.
	 */
	logErr("(%d)%s: exec failed (%s)\n", err, err, cmd)
	os.Exit(255)
}
