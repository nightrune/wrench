package main

import (
  "fmt"
  "github.com/nightrune/wrench/logging"
  "flag"
  "os"
  "strings"
  "encoding/json"
  "io/ioutil"
)

/*
 Example build file
 {
    agent_file:"path/to/file"
    device_file:"path/to/file"
    dirs:[
      "libs"
    ]
 }
 */
const DEFAULT_PROJECT_FILE = "settings.wrench"
const DEFAULT_LIB_DIR = "libs"
const DEFAULT_DEVICE_IN_FILE = "device.nut.in"
const DEFAULT_AGENT_IN_FILE = "agent.nut.in"
const DEFAULT_DEVICE_OUT_FILE = "device.nut"
const DEFAULT_AGENT_OUT_FILE = "agent.nut"

type ProjectSettings struct {
  AgentFileOutPath string `json:"agent_file_out"`
  DeviceFileOutPath string `json:"device_file_out"`
  AgentFileInPath string `json:"agent_file_in"`
  DeviceFileInPath string `json:"device_file_in"`
  ApiKeyFile string `json:"api_file"`
  ModelKey string `json:"model_key"`
  LibraryDirs []string `json:"dirs"`
}
 
// A Command is an implementation of a go command
// like go build or go fix.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool
  
  settings ProjectSettings
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "  %s\n", strings.TrimSpace(c.Long))
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var commands = []*Command {
  cmdBuild,
  cmdRun,
  cmdTest,
  cmdUpload,
  cmdModel,
}

func PrintHelp() {
  fmt.Printf("Help: \n")
  flag.PrintDefaults()
  for _, cmd := range commands {
    if cmd.Runnable() {
      cmd.Usage()
    }
  }
}

func LoadSettings(settings_file string) ProjectSettings {
  settings_data, err := ioutil.ReadFile(settings_file)
  if err != nil {
    logging.Fatal("Could not open settings file %s", settings_file)
    os.Exit(1)
  }
  
  var settings ProjectSettings
  err = json.Unmarshal(settings_data, &settings)
  if err != nil {
    logging.Fatal("Couldn't parse settings file: %s", settings_file)
    os.Exit(1)
  }
  return settings
}

func ProcessSettings(settings_file string) ProjectSettings {
  if _, err := os.Stat(settings_file); err != nil {
    logging.Info("Did not find the settings file %s", settings_file)
    if settings_file != DEFAULT_PROJECT_FILE {
      logging.Fatal("Could not load non default settings file...")
      os.Exit(1)
    }
    logging.Info("Generating default settings file...")
    var settings ProjectSettings
    settings.AgentFileInPath = DEFAULT_AGENT_IN_FILE
    settings.AgentFileOutPath = DEFAULT_AGENT_OUT_FILE
    settings.DeviceFileInPath = DEFAULT_DEVICE_IN_FILE
    settings.DeviceFileOutPath = DEFAULT_DEVICE_OUT_FILE
    settings.LibraryDirs = append(settings.LibraryDirs, DEFAULT_LIB_DIR)
    settings_data, err := json.Marshal(settings)
    if err != nil {
      logging.Warn("Failed to generate default settings json data")
    } else {
      err = ioutil.WriteFile(settings_file, settings_data, 777)
      if err != nil {
        logging.Warn("Failed to write new defaults...")
      }
    }
    
    return settings
  }
  
  return LoadSettings(settings_file)
}

func main() {
  logging_int := flag.Int("l", int(logging.LOG_INFO),
    "Levels 0-4 0 == None, 4 == Debug")
  log_colors_flag := flag.Bool("log_color", false, "-log_color enables log coloring(mingw/linux only)")
  settings_file := flag.String("settings", DEFAULT_PROJECT_FILE,
    "Set the settings file to a non standard file...")

  flag.Parse()
  err, log_value := logging.IntToLogLevel(*logging_int);
  if err == nil {
    logging.SetLoggingLevel(log_value);
  } else {
    PrintHelp()
    os.Exit(1)
    return
  }

  logging.SetColorEnabled(*log_colors_flag)

  args := flag.Args()
  
  if len(args) < 1 {
    logging.Info("Need a subcommand")
    PrintHelp()
    os.Exit(1)
    return
  }
  
  logging.Debug("Using settings file: %s", *settings_file)
  projectSettings := ProcessSettings(*settings_file)
  logging.Debug("Settings found: %s", projectSettings)
  for i, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[i:]
			} else if len(args) > 2 {
				cmd.Flag.Parse(args[i:])
				args = cmd.Flag.Args()
			}
      cmd.settings = projectSettings;
			cmd.Run(cmd, args)
			return
		}
	}
}
