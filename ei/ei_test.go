package ei

import "github.com/nightrune/wrench/logging"
import "fmt"

// These are broken because you need to be able to load values
const API_KEY = ""
const MODEL_KEY = ""
const TEST_DEVICE_ID = ""

func ExampleGetModel() {
	logging.SetLoggingLevel(logging.LOG_DEBUG)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	model_list, err := client.ListModels()
	if err != nil {
		logging.Fatal("Test Failed %s", err.Error())
		return
	}

	for _, model := range model_list.Models {
		logging.Info("Id: %s, Name: %s", model.Id, model.Name)
	}
	logging.Info("Ending Test")

}

func ExampleCreateUpdateDeleteModel() {
	logging.SetLoggingLevel(logging.LOG_DEBUG)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	new_model := new(Model)
	new_model.Name = "Wrench Test Model"
	model, err := client.CreateModel(new_model)
	if err != nil {
		logging.Fatal("Example Failed %s", err.Error())
		return
	}

	new_model.Name = "Wrench Test Model Again"
	update_model, err := client.UpdateModel(model.Id, new_model)
	if err != nil {
		logging.Fatal("Example Failed %s", err.Error())
		return
	}

	if update_model.Id != model.Id {
		logging.Fatal("Example Failed, Model Ids don't match")
		return
	}

	err = client.DeleteModel(model.Id)
	if err != nil {
		logging.Fatal("Example Failed %s", err.Error())
		return
	}
	logging.Info("Succesfully created model: %s with Id: %s", model.Name, model.Id)
	logging.Info("Ending Test")

}

func ExampleRestartModelDevices() {
	logging.SetLoggingLevel(logging.LOG_DEBUG)
	logging.Info("Starting Example")
	client := NewBuildClient(API_KEY)
	err := client.RestartModelDevices(MODEL_KEY)
	if err != nil {
		logging.Fatal("Example Failed %s", err.Error())
		return
	}

	logging.Info("Succesfully restarted model: %s", MODEL_KEY)
	logging.Info("Ending Example")

}

func ExampleGetCodeRevisions() {
	logging.SetLoggingLevel(logging.LOG_DEBUG)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	revisions, err := client.GetCodeRevisionList(MODEL_KEY)
	if err != nil {
		logging.Fatal("Test Failed %s", err.Error())
		return
	}

	for _, revision := range revisions {
		logging.Info(string(revision.CreatedAt))
		logging.Info(`Version: %d
  	              CreatedAt: %s,
	  	          ReleaseNotes: %s,`,
			revision.Version,
			revision.CreatedAt,
			revision.ReleaseNotes)
	}
	logging.Info("Ending Test")

}

func ExampleGetCodeRevision() {
	logging.SetLoggingLevel(logging.LOG_DEBUG)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	revision, err := client.GetCodeRevision(MODEL_KEY, "1")
	if err != nil {
		logging.Fatal("Test Failed %s", err.Error())
		return
	}

	logging.Info(string(revision.CreatedAt))
	logging.Info(`Version: %d
  	            CreatedAt: %s,
  	            ReleaseNotes: %s,
  	            AgentCode: %s,
  	            DeviceCode: %s,`,
		revision.Version,
		revision.CreatedAt,
		revision.ReleaseNotes,
		revision.AgentCode,
		revision.DeviceCode)
	logging.Info("Ending Test")

}

func ExampleUploadCode() {
	logging.SetLoggingLevel(logging.LOG_DEBUG)
	logging.Info("Starting Test")
	request := new(CodeRevisionLong)
	client := NewBuildClient(API_KEY)
	request.AgentCode = `server.log("More Agent Code!")`
	request.DeviceCode = `server.log("More Device Code!")`
	client.UpdateCodeRevision(MODEL_KEY, request)
	logging.Info("Ending Test")

}

func ExampleDeviceList() {
	logging.SetLoggingLevel(logging.LOG_INFO)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	list, err := client.GetDeviceList()
	if err != nil {
		logging.Fatal("Failed to get device list! %s", err.Error())
		return
	}
	for _, device := range list {
		logging.Info("Name: %s, Id: %s, Model: %s", device.Name, device.Id, device.ModelId)
	}
	logging.Info("Ending Test")

}

func ExampleGetDeviceLogs() {
	logging.SetLoggingLevel(logging.LOG_INFO)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	logs, err := client.GetDeviceLogs(TEST_DEVICE_ID)
	if err != nil {
		logging.Fatal("Failed to get device list! %s", err.Error())
		return
	}
	for _, log := range logs {
		fmt.Printf("%s %s:%s\n", log.Timestamp, log.Type, log.Message)
	}
	logging.Info("Ending Test")

}

func ExampleGetDevice() {
	logging.SetLoggingLevel(logging.LOG_INFO)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	device, err := client.GetDevice(TEST_DEVICE_ID)
	if err != nil {
		logging.Fatal("Failed to get device! %s", err.Error())
		return
	}
	fmt.Printf("Name: %s Id: %s\n", device.Name, device.Id)
	logging.Info("Ending Test")
}

func ExampleRestartDevice() {
	logging.SetLoggingLevel(logging.LOG_INFO)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	err := client.RestartDevice(TEST_DEVICE_ID)
	if err != nil {
		logging.Fatal("Failed to get device! %s", err.Error())
		return
	}
	fmt.Printf("Device Reset")
	logging.Info("Ending Test")
}

func ExampleUpdateDevice() {
	logging.SetLoggingLevel(logging.LOG_INFO)
	logging.Info("Starting Test")
	client := NewBuildClient(API_KEY)
	new_device := new(Device)
	new_device.Name = "Wiggy Whacky - hamburg-2"
	dev, err := client.UpdateDevice(new_device, TEST_DEVICE_ID)
	if err != nil {
		logging.Fatal("Failed to update device! %s", err.Error())
		return
	}
	fmt.Printf("Device Updated: %s", dev.Name)
	logging.Info("Ending Test")
}
