package main

import (
	"github.com/hiboma/suexec"
	"os"
)

func main() {

	environ := suexec.CleanEnv()

	log := suexec.NewLog(suexec.AP_LOG_EXEC)

	cwd, err := os.Getwd()
	if err != nil {
		log.LogErr("cannot get current working directory\n")
		os.Exit(111)
	}

	p := suexec.Param{
		Args: os.Args,
		Uid:  os.Getuid(),
		Cwd:  cwd,
		Log:  log,
	}

	script := suexec.NewSuexec(p)

	if err := script.VerifyToSuexec(); err != nil {
		if err.Status() != 0 {
			log.LogErr("%s\n", err.Message())
		}
		os.Exit(err.Status())
	}

	if err := script.Exec(environ); err != nil {
		log.LogErr("%s\n", err.Message())
		os.Exit(err.Status())
	}

	os.Exit(255)
}
