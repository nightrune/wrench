package ei

import "github.com/nightrune/wrench/logging"

const api_key = "1d2da2b7e4e35667283af41ba2458527"
const model_key = "VqNg23hFUNjg"

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