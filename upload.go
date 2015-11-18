package main

import (
	"encoding/json"
	"github.com/nightrune/wrench/ei"
	"github.com/nightrune/wrench/logging"
	"io/ioutil"
	"os"
)

type ApiKeyFile struct {
	Key string `json:"key"`
}

var cmdUpload = &Command{
	UsageLine: "upload",
	Short:     "Upload files with api key and set model",
	Long:      "Uploads the files agent.nut and device.nut into the model selected within settings.wrench",
}

func init() {
	cmdUpload.Run = UploadFiles
}

func UploadFiles(cmd *Command, args []string) {
	logging.Info("Attempting to upload...")
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

	agent_code, err := ioutil.ReadFile(cmd.settings.AgentFileOutPath)
	if err != nil {
		logging.Fatal("Could not open the agent code: %s", cmd.settings.AgentFileOutPath)
		os.Exit(1)
		return
	}

	device_code, err := ioutil.ReadFile(cmd.settings.DeviceFileOutPath)
	if err != nil {
		logging.Fatal("Could not open the device code %s", cmd.settings.DeviceFileOutPath)
		os.Exit(1)
		return
	}

	request := new(ei.CodeRevisionLong)
	client := ei.NewBuildClient(keyfile.Key)
	request.AgentCode = string(agent_code)
	request.DeviceCode = string(device_code)
	response, err := client.UpdateCodeRevision(cmd.settings.ModelKey, request)
	if err != nil {
		logging.Fatal("Failed to upload code to model %s, Error: %s", cmd.settings.ModelKey, err.Error())
		os.Exit(1)
		return
	}
	logging.Info("Succesfully uploaded version %d", response.Version)
}
