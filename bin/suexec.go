package main

import (
	"github.com/hiboma/suexec"
	"os"
)

func main() {
	status := suexec.Suexec(os.Args, suexec.NewLog(suexec.AP_LOG_EXEC))
	os.Exit(status)
}
