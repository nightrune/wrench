package main

import (
	"encoding/json"
	"fmt"
	"github.com/nightrune/wrench/ei"
	"github.com/nightrune/wrench/logging"
	"io/ioutil"
	"os"
)

var cmdModel = &Command{
	UsageLine: "model",
	Short:     "Various ways to interact with electric imp models",
	Long: `Submcommand:
    model - Shows usage and current model name
    model list - Lists the current models`,
}

func init() {
	cmdModel.Run = ModelSubCommand
}

func ListModels(api_key string) {
	client := ei.NewBuildClient(api_key)
	model_list, err := client.ListModels()
	if err != nil {
		logging.Fatal("Failed to get model list %s", err.Error())
		return
	}

	for _, model := range model_list.Models {
		fmt.Printf("Id: %s, Name: %s\n", model.Id, model.Name)
	}
}

// TODO(sean) Rewrite this, and make the command list from wrench.go portable and resulable
func PrintModelHelp() {
	fmt.Printf("Available sub-commands, list\n")
}

func ModelSubCommand(cmd *Command, args []string) {
	logging.Debug("In model")
	for _, s := range args {
		logging.Debug(s)
	}

	// TODO(sean) Break this out to a helpers area for the ei stuff
	keyfile_data, err := ioutil.ReadFile(cmd.settings.ApiKeyFile)
	if err != nil {
		logging.Fatal("Could not open the keyfile: %s", cmd.settings.ApiKeyFile)
		os.Exit(1)
		return
	}

	keyfile := new(ApiKeyFile)
	err = json.Unmarshal(keyfile_data, keyfile)
	if err != nil {
		logging.Fatal("Could not parse keyfile: %s", cmd.settings.ApiKeyFile)
		os.Exit(1)
		return
	}

	client := ei.NewBuildClient(keyfile.Key)

	if len(args) < 2 {
		logging.Debug(cmd.settings.ModelKey)
		if cmd.settings.ModelKey != "" {
			model, err := client.GetModel(cmd.settings.ModelKey)
			if err != nil {
				logging.Fatal("Failed to get model id: %s, Got Error: %s",
					cmd.settings.ModelKey, err.Error())
				os.Exit(1)
				return
			}
			fmt.Printf("Current Model: %s\n", model.Name)
		} else {
			fmt.Printf("No Model set in settings.wrench\n")
		}

		PrintModelHelp()
		os.Exit(1)
	}

	if args[1] == "list" {
		ListModels(keyfile.Key)
	} else {
		PrintModelHelp()
		os.Exit(1)
	}
}
