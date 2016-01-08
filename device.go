package main

import (
	"fmt"
	"github.com/nightrune/wrench/ei"
	"github.com/nightrune/wrench/logging"
	"os"
)

var cmdDevice = &Command{
	UsageLine: "device",
	Short:     "Various ways to interact with electric imp devices",
	Long: `Submcommand:
    device - Shows usage
    device list - Lists the current models
    device assign <device_id> <model_id>
    device restart <device_id>`,
}

func init() {
	cmdDevice.Run = DeviceSubCommand
}

func ListDevices(client *ei.BuildClient) {
	device_list, err := client.GetDeviceList()
	if err != nil {
		logging.Fatal("Failed to get device list %s", err.Error())
		return
	}

	for _, model := range device_list {
		fmt.Printf("Id: %s, Name: %s\n", model.Id, model.Name)
	}
}

func StartDeviceLogging(client *ei.BuildClient, device_id string) {
	logs, poll_url, err := client.GetDeviceLogs(device_id)
	if err != nil {
		logging.Fatal("Failed to get device logs %s", err.Error())
		return
	}

    logging.Debug("Poll Url %s", poll_url)
	for _, entry := range logs {
		fmt.Printf("%s: %s: %s\n", entry.Timestamp, entry.Type, entry.Message)
	}

	for {
		logs, err = client.ContinueDeviceLogs(poll_url)
		if err != nil {
			if _, ok := err.(*ei.Timeout); ok == false {
				logging.Fatal("Failed to get device logs %s", err.Error())
				return;
			} else {
				logging.Debug("Long poll timed out...")
			}
		} else {
			for _, entry := range logs {
				fmt.Printf("%s: %s: %s\n", entry.Timestamp, entry.Type, entry.Message)
			}
		}
	}
}

func AssignDeviceToModel(client *ei.BuildClient, device_id string, model_id string) {
	device, err := client.GetDevice(device_id)
	if err != nil {
		logging.Fatal("Failed to retrieve current model for model id: %s", model_id)
		return
	}

	// Don't do anything if its already the same
	if device.ModelId == model_id {
		return;
	}

	device.ModelId = model_id;
	updated_device, err := client.UpdateDevice(&device, device_id)
	logging.Debug("Model and new Devices: %s", updated_device)
}


func RestartDevice(client *ei.BuildClient, device_id string) {
	err := client.RestartDevice(device_id)
	if err != nil {
		logging.Fatal("Failed to restart device %s with error %s", device_id, err.Error())
		return
	}
	fmt.Printf("Successfully restarted device and agent");
}

// TODO(sean) Rewrite this, and make the command list from wrench.go portable and resulable
func PrintDeviceHelp() {
	fmt.Printf("Available sub-commands, list, log, assign, restart\n")
}

func DeviceSubCommand(cmd *Command, args []string) {
	logging.Debug("In Device")
	for i, s := range args {
		logging.Debug("%d:%s", i, s)
	}

	client, err := CreateClient(cmd.settings.ApiKeyFile)
	if err != nil {
		os.Exit(1)
		return
	}

	if len(args) < 2 {
		PrintDeviceHelp()
		os.Exit(1)
	}

	if args[1] == "list" {
		ListDevices(client)
	} else if args[1] == "log" {
		if len(args) < 3 {
			fmt.Printf("Usage: device log <device_id>\n")
			PrintDeviceHelp()
		} else {
			logging.Debug("Attempting to start logging")
			StartDeviceLogging(client, args[2])
		}
	} else if args[1] == "assign" {
		if len(args) < 4 {
			fmt.Printf("Usage: device assign <device_id> <model_id>\n");
			PrintDeviceHelp();
		} else {
			logging.Debug("Attempting to assign device to model");
			AssignDeviceToModel(client, args[2], args[3]);
		}
	} else if args[1] == "restart" {
		if len(args) < 3 {
			fmt.Printf("Usage: device restart <device_id>\n");
			PrintDeviceHelp();
		} else {
			logging.Debug("Attempting to restart device");
			RestartDevice(client, args[2]);
		}
	} else {
		logging.Debug("Showing help")
		PrintDeviceHelp()
		os.Exit(1)
	}
}
