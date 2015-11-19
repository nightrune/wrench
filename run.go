package main

import (
	"github.com/nightrune/wrench/logging"
	"github.com/nightrune/wrench/squirrel"
	"os"
)

const SQRL_CMD_BIN = "sq.exe"

var cmdRun = &Command{
	UsageLine: "run",
	Short:     "Runs a file",
	Long:      "run [input_file]",
}

func init() {
	cmdRun.Run = RunMe
}

func ExecuteSqrl(script string) error {
	squirrel.RunScript(script)
	return nil
}

func RunMe(cmd *Command, args []string) {
	logging.Info("Attempting to run script...")
	if len(args) < 2 {
		logging.Fatal("run requires file name")
		os.Exit(1)
	}
	ExecuteSqrl(args[1])
	logging.Info("Script finished")
}
