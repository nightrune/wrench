package ei

import "github.com/nightrune/wrench/logging"
import "fmt"
const api_key = "1d2da2b7e4e35667283af41ba2458527"
const model_key = "VqNg23hFUNjg"
const model_key_delete = "OqtySgkS-UD8"

func ExampleGetModel() {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Test")
  client := NewBuildClient(api_key)
  model_list, err := client.ListModels();
  if err != nil {
  	logging.Fatal("Test Failed %s", err.Error())
  	return
  }

  for _, model := range model_list.Models {
  	logging.Info("Id: %s, Name: %s", model.Id, model.Name)
  }
  logging.Info("Ending Test")
  //Output:
}

func ExampleCreateModel() {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Test")
  client := NewBuildClient(api_key)
  new_model := new(Model)
  new_model.Name = "Wrench Test Model"
  model, err := client.CreateModel(new_model);
  if err != nil {
    logging.Fatal("Example Failed %s", err.Error())
    return
  }

  logging.Info("Succesfully created model: %s with Id: %s", model.Name, model.Id)
  logging.Info("Ending Test")
  //Output:
}

func ExampleCreateUpdateDeleteModel() {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Test")
  client := NewBuildClient(api_key)
  new_model := new(Model)
  new_model.Name = "Wrench Test Model"
  model, err := client.CreateModel(new_model);
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
  //Output:
}

func ExampleDeleteModel() {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Example")
  client := NewBuildClient(api_key)
  err := client.DeleteModel(model_key_delete);
  if err != nil {
    logging.Fatal("Example Failed %s", err.Error())
    return
  }

  logging.Info("Succesfully deleted model: %s", model_key_delete)
  logging.Info("Ending Example")
  //Output:
}

func ExampleRestartModelDevices() {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Example")
  client := NewBuildClient(api_key)
  err := client.RestartModelDevices(model_key);
  if err != nil {
    logging.Fatal("Example Failed %s", err.Error())
    return
  }

  logging.Info("Succesfully restarted model: %s", model_key)
  logging.Info("Ending Example")
  //Output:
}

func ExampleGetCodeRevisions() {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Test")
  client := NewBuildClient(api_key)
  revisions, err := client.GetCodeRevisionList(model_key);
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
  //Output:
}

func ExampleGetCodeRevision() {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Test")
  client := NewBuildClient(api_key)
  revision, err := client.GetCodeRevision(model_key, "1");
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
  //Output:
}

func ExampleUploadCode() {
	logging.SetLoggingLevel(logging.LOG_DEBUG);
  	logging.Info("Starting Test")
	request := new(CodeRevisionLong)
	client := NewBuildClient(api_key)
	request.AgentCode = `server.log("More Agent Code!")`
	request.DeviceCode = `server.log("More Device Code!")`
  	client.UpdateCodeRevision(model_key, request);
	logging.Info("Ending Test")
	//Output:
}

func ExampleDeviceList() {
	logging.SetLoggingLevel(logging.LOG_INFO);
  	logging.Info("Starting Test")
	client := NewBuildClient(api_key)
  	list, err := client.GetDeviceList()
  	if err != nil {
  		logging.Fatal("Failed to get device list! %s", err.Error())
  		return;
  	}
  	for _, device := range list {
  		logging.Info("Name: %s, Id: %s, Model: %s", device.Name, device.Id, device.ModelId)
  	}
	logging.Info("Ending Test")
	//Output:
}

const TEST_DEVICE_ID = "30000c2a690be1e1"
func ExampleGetDeviceLogs() {
	logging.SetLoggingLevel(logging.LOG_INFO);
  	logging.Info("Starting Test")
	client := NewBuildClient(api_key)
  	logs, err := client.GetDeviceLogs(TEST_DEVICE_ID)
  	if err != nil {
  		logging.Fatal("Failed to get device list! %s", err.Error())
  		return;
  	}
  	for _, log := range logs {
  		fmt.Printf("%s %s:%s\n", log.Timestamp, log.Type, log.Message)
  	}
	logging.Info("Ending Test")
	//Output:
}