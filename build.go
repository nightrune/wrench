package main

import (
  "fmt"
  "github.com/gyuho/goraph"
)

var cmdBuild = &Command {
  UsageLine: "build",
  Short: "Preprocess Squirel",
  Long:"",
}

func init() {
  cmdBuild.Run = RunMe
}

func ProcessLineMacro()

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
  
  // Start with the agent file, and device files
  
  // Get a list of all the #import statements
    // Then open and scan all imported files for a #library name
    //   If a file is invalid throw error
    // Make a graph of library and file imports
  
  // Check for acyclical graph with a topological sort
  // If loops fail
  // Else keep going
  
  // Now that we have a list of includes
  // Preprocess them for __FILE__, and __LINE__
  
  
}