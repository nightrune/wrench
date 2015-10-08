package main

import (
  "fmt"
)

var cmdBuild = &Command {
  UsageLine: "build",
  Short: "Preprocess Squirel",
  Long:"",
}

func init() {
  cmdBuild.Run = RunMe
}

func RunMe(cmd *Command, args []string) {
  fmt.Printf("BUILDING!!\n")
}