package main

import (
  "os/exec"
  "os"
  "bufio"
  "github.com/nightrune/wrench/logging"
  "errors"
  "strings"
)

const PREPROCESSOR_CMD_BIN = "gpp"

var cmdBuild = &Command {
  UsageLine: "build",
  Short: "Preprocess Squirel",
  Long:"",
}

func init() {
  cmdBuild.Run = BuildMe
}

func PreProcessFile(outputFile string, inputFile string, libraryDirs []string) error {
  if _, err := os.Stat(inputFile); err != nil {
      logging.Fatal("Did not find input file %s", inputFile)
      return errors.New("Input file not found")
  }
  gppCmdName := PREPROCESSOR_CMD_BIN
  gppArgs := []string{"-o", outputFile, inputFile}
  var s []string
  for _, dir := range libraryDirs {
    s = []string{"-I", dir}
    gppArgs = append(gppArgs, strings.Join(s, ""))
  }
  gppArgs = append(gppArgs, "-C")
  logging.Debug("gppArgs: %s", gppArgs)
  gppCmd := exec.Command(gppCmdName, gppArgs...)
  cmdReader, err := gppCmd.StdoutPipe()
  if err != nil {
    logging.Fatal("Failed to get stdout pipe Error: %s\n", err.Error())
    os.Exit(1);
  }
  
  cmdErrorReader, err := gppCmd.StderrPipe()
  if err != nil {
    logging.Fatal("Failed to get stderr pipe Error: %s\n", err.Error())
    os.Exit(1);
  }
  
  scanner := bufio.NewScanner(cmdReader)
  go func() {
    for scanner.Scan() {
      logging.Info("%s", scanner.Text())
    }
  }()
  
  errScanner := bufio.NewScanner(cmdErrorReader)
  go func() {
    for errScanner.Scan() {
      logging.Warn("%s", errScanner.Text())
    }
  }()
  
  err = gppCmd.Start()
	if err != nil {
		logging.Fatal("Error starting gpp", err)
		os.Exit(1)
	}

	err = gppCmd.Wait()
	if err != nil {
		logging.Fatal("Error waiting for gpp", err)
		os.Exit(1)
	}
  
  return nil;
}
  
func BuildMe(cmd *Command, args []string) {
  logging.Info("Starting Build Process...")
  err := PreProcessFile(cmd.settings.AgentFileOutPath,
    cmd.settings.AgentFileInPath, cmd.settings.LibraryDirs);
  if err != nil {
    os.Exit(1);
  }
  
  err = PreProcessFile(cmd.settings.DeviceFileOutPath,
    cmd.settings.DeviceFileInPath, cmd.settings.LibraryDirs);
  if err != nil {
    os.Exit(1);
  }
  logging.Info("Succesfully built...")
}
