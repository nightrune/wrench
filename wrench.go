package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nightrune/wrench/logging"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
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
const DEFAULT_DEVICE_IN_FILE = "device.nut"
const DEFAULT_AGENT_IN_FILE = "agent.nut"
const DEFAULT_DEVICE_OUT_FILE = "device.nut.out"
const DEFAULT_AGENT_OUT_FILE = "agent.nut.out"
const DEFAULT_API_KEY_FILE = ".build_api_key.json"

type ProjectSettings struct {
	AgentFileOutPath  string   `json:"agent_file_out"`
	DeviceFileOutPath string   `json:"device_file_out"`
	AgentFileInPath   string   `json:"agent_file_in"`
	DeviceFileInPath  string   `json:"device_file_in"`
	ApiKeyFile        string   `json:"api_file"`
	ModelKey          string   `json:"model_key"`
	LibraryDirs       []string `json:"dirs"`
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
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "  %s\n\n", strings.TrimSpace(c.Long))
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var commands = []*Command{
	cmdBuild,
	cmdRun,
	cmdTest,
	cmdUpload,
	cmdModel,
	cmdDevice,
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

func LoadSettings(settings_file string) (ProjectSettings, error) {
	var settings ProjectSettings
	settings_data, err := ioutil.ReadFile(settings_file)
	if err != nil {
		logging.Fatal("Could not open settings file %s", settings_file)
		return settings, err
	}

	err = json.Unmarshal(settings_data, &settings)
	if err != nil {
		logging.Fatal("Couldn't parse settings file: %s", settings_file)
		return settings, err
	}
	return settings, nil
}

func GenerateDefaultSettingsFile(settings_file string) (ProjectSettings, error) {
	logging.Info("Generating default settings file...")
	var settings ProjectSettings
	settings.AgentFileInPath = DEFAULT_AGENT_IN_FILE
	settings.AgentFileOutPath = DEFAULT_AGENT_OUT_FILE
	settings.DeviceFileInPath = DEFAULT_DEVICE_IN_FILE
	settings.DeviceFileOutPath = DEFAULT_DEVICE_OUT_FILE
	settings.ApiKeyFile = DEFAULT_API_KEY_FILE
	settings.LibraryDirs = append(settings.LibraryDirs, DEFAULT_LIB_DIR)
	settings_data, err := json.Marshal(settings)
	if err != nil {
		logging.Warn("Failed to generate default settings json data")
		return settings, err
	} else {
		err = ioutil.WriteFile(settings_file, settings_data, 777)
		if err != nil {
			logging.Warn("Failed to write new defaults...")
			return settings, err
		}
	}
	return settings, nil
}

func main() {
	logging_int := flag.Int("l", int(logging.LOG_INFO),
		"Levels 0-4 0 == None, 4 == Debug")
	log_colors_flag := flag.Bool("log_color", false, "-log_color enables log coloring(mingw/linux only)")
	settings_file := flag.String("settings", DEFAULT_PROJECT_FILE,
		"Set the settings file to a non standard file...")
	gen_settings_flag := flag.Bool("g", false, "-g generates a default settings file if one is not found")

	flag.Parse()
	err, log_value := logging.IntToLogLevel(*logging_int)
	if err == nil {
		logging.SetLoggingLevel(log_value)
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
	    for {
	    	select {
	    	case <-c:
	    		os.Exit(0)
	    	}
	    }
	}()

    var projectSettings ProjectSettings
	logging.Debug("Using settings file: %s", *settings_file)
	if _, err := os.Stat(*settings_file); err != nil {
		logging.Info("Did not find the settings file %s", *settings_file)
		if *settings_file != DEFAULT_PROJECT_FILE {
			logging.Fatal("Could not load non default settings file...")
			os.Exit(1)
		}

		if (*gen_settings_flag) {
			projectSettings, err = GenerateDefaultSettingsFile(*settings_file)
			if (err != nil) {
				os.Exit(1)
			}
		} else {
			os.Exit(1);
		}
	} else {
		projectSettings, err = LoadSettings(*settings_file);
		if (err != nil) {
			logging.Info("Failed to load settings file: %s", *settings_file)
			os.Exit(1)
		}
	}

	logging.Debug("Settings found: %s", projectSettings)
	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmd.Flag.Usage = func() { cmd.Usage() }
			for i, s := range args {
				logging.Debug("Left Args: %d:%s", i, s)
			}
			if cmd.CustomFlags {
				args = args[0:]
			} else if len(args) > 2 {
				cmd.Flag.Parse(args[0:])
				args = cmd.Flag.Args()
			}
			cmd.settings = projectSettings
			cmd.Run(cmd, args)
			return
		}
	}
}
