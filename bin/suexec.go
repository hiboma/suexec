package main

import (
	"github.com/hiboma/suexec"
	"os"
)

func main() {

	log := suexec.NewLog(suexec.AP_LOG_EXEC)

	cwd, err := os.Getwd()
	if err != nil {
		log.LogErr("cannot get current working directory\n")
		os.Exit(111)
	}

	status := suexec.Suexec(suexec.Param{args: os.Args, uid: os.Getuid, cwd: cwd, log: log})
	os.Exit(status)
}
