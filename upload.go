package main

import (
	"encoding/json"
	"github.com/nightrune/wrench/ei"
	"github.com/nightrune/wrench/logging"
	"io/ioutil"
	"os"
	"fmt"
	"flag"
)

type ApiKeyFile struct {
	Key string `json:"key"`
}

var cmdUpload = &Command{
	UsageLine: "upload <-r>",
	Short:     "Upload files with api key and set model",
	Long:      `Uploads the files agent.nut and device.nut into the model selected within settings.wrench
	            -r Allows you to restart the device after an upload, Restarts all devices on model`,
}

func init() {
	cmdUpload.Run = UploadFiles
}

func UploadUsage() {

}

func UploadFiles(cmd *Command, args []string) {
	// Create our flags
	flag_set := flag.NewFlagSet("UploadFlagSet", flag.ExitOnError)
	restart_device := flag_set.Bool("r", false, "-r restarts the server on a successful model upload");
	flag_set.Parse(args[1:]);

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

	if response.Success == false {
		fmt.Println("There were errors in the build:")
		fmt.Printf("Error Code: %s \n", response.Error.Code);
		fmt.Printf("Message: %s \n", response.Error.FullMessage);
		return
	}

	logging.Info("Succesfully uploaded version %d", response.Revisions.Version)

	if *restart_device == true {
		err := client.RestartModelDevices(cmd.settings.ModelKey);
		if err != nil {
			logging.Fatal("Failed to restart devices after upload, Error: %s", err.Error());
			os.Exit(1)
		}

		model, err := client.GetModel(cmd.settings.ModelKey)
		if err != nil {
			fmt.Printf("Model: %s devices restarted.\n", cmd.settings.ModelKey);
			return
		} else {
			fmt.Printf("Model: %s restarted.\n", model.Name)
		}
	}
}
