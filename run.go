package main

import (
  "fmt"
  "os/exec"
  "os"
  "bufio"
  "github.com/nightrune/wrench/logging"
)

const SQRL_CMD_BIN = "./sq.exe"

var cmdRun = &Command {
  UsageLine: "run",
  Short: "Runs a file",
  Long:"run [input_file]",
}

func init() {
  cmdRun.Run = RunMe
}

func ExecuteSqrl(script string) error {
  sqrlCmdName := SQRL_CMD_BIN
  sqrlArgs := []string{script}
  sqrlCmd := exec.Command(sqrlCmdName, sqrlArgs...)
  cmdReader, err := sqrlCmd.StdoutPipe()
  if err != nil {
    logging.Fatal("Failed to get stdout pipe Error: %s\n", err.Error())
    os.Exit(1);
  }
  
  sqrlErrorReader, err := sqrlCmd.StderrPipe()
  if err != nil {
    logging.Fatal("Failed to get stderr pipe Error: %s", err.Error())
    os.Exit(1);
  }
  
  scanner := bufio.NewScanner(cmdReader)
  go func() {
    for scanner.Scan() {
      fmt.Printf("--%q\n", scanner.Text())
    }
  }()
  
  errScanner := bufio.NewScanner(sqrlErrorReader)
  go func() {
    for errScanner.Scan() {
      fmt.Printf("--%q\n", errScanner.Text())
    }
  }()
  
  err = sqrlCmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting squirrel", err)
		os.Exit(1)
	}

	err = sqrlCmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for squirrel", err)
		os.Exit(1)
	}
  
  return nil;
}

func RunMe(cmd *Command, args []string) {
  logging.Info("Attempting to run script...")
  if len(args) < 1 {
    logging.Fatal("run requires file name")
    os.Exit(1)
  }
  ExecuteSqrl(args[0])
  logging.Info("Script finished")
}