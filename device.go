package main

import (
	"fmt"
	"github.com/nightrune/wrench/ei"
	"github.com/nightrune/wrench/logging"
	"os"
	"reflect"
)

var cmdDevice = &Command{
	UsageLine: "device",
	Short:     "Various ways to interact with electric imp devices",
	Long: `Submcommand:
    device - Shows usage
    device list - Lists the current models`,
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

// TODO(sean) Rewrite this, and make the command list from wrench.go portable and resulable
func PrintDeviceHelp() {
	fmt.Printf("Available sub-commands, list\n")
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
		PrintModelHelp()
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
	} else {
		logging.Debug("Showing help")
		PrintDeviceHelp()
		os.Exit(1)
	}
}
