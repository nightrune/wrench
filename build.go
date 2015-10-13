package main

import (
  "fmt"
  "os/exec"
  "os"
  "bufio"
)

var cmdBuild = &Command {
  UsageLine: "build",
  Short: "Preprocess Squirel",
  Long:"",
}

func init() {
  cmdBuild.Run = RunMe
}

func ProcessLineMacro() {
}

/*
 Example build file
 {
    agent_file:"path/to/file"
    device_file:"path/to/file"
    library_dirs:[
      "libs"
    ]
 }
 */

func RunMe(cmd *Command, args []string) {
  fmt.Printf("BUILDING!!\n")
  // Load the build file
  
  // Get a list of all the #import statements
    // Then open and scan all imported files for a #library name
    //   If a file is invalid throw error
    // Make a graph of library and file imports
  
  // Check for acyclical graph with a topological sort
  // If loops fail
  // Else keep going
  
  // Now that we have a list of includes
  // Preprocess them for __FILE__, and __LINE__
  
  gppCmdName := "./gpp"
  gppArgs := []string{"-o", "test.out.nut", "./test_sql/test.nut", "-I./test_sql"}
  //gppArgs := []string{"-h"}
  gppCmd := exec.Command(gppCmdName, gppArgs...)
  cmdReader, err := gppCmd.StdoutPipe()
  if err != nil {
    fmt.Printf("Failed to get stdout pipe Error: %s\n", err.Error())
    os.Exit(1);
  }
  
  cmdErrorReader, err := gppCmd.StderrPipe()
  if err != nil {
    fmt.Printf("Failed to get stderr pipe Error: %s\n", err.Error())
    os.Exit(1);
  }
  
  scanner := bufio.NewScanner(cmdReader)
  go func() {
    for scanner.Scan() {
      fmt.Printf("%s\n", scanner.Text())
    }
  }()
  
  errScanner := bufio.NewScanner(cmdErrorReader)
  go func() {
    for errScanner.Scan() {
      fmt.Printf("%s\n", errScanner.Text())
    }
  }()
  
  err = gppCmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting gpp", err)
		os.Exit(1)
	}

	err = gppCmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for gpp", err)
		os.Exit(1)
    return;
	}
  
  fmt.Printf("Succesfully built?")
}