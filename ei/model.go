package ei

import (
	"net/http"
	"bytes"
	"github.com/nightrune/wrench/logging"
)

const EI_URL = "https://build.electricimp.com/v4/"

type Device struct {
  Id string `json:"id"`
  Name string `json:"name"`
  ModelId string `json:"model_id"`
  PowerState string `json:"powerstate"`
  Rssi int `json:"rssi"`
  AgentId string `json:"agent_id"`
  AgentStatus string `json:"agent_status"`
}

type Model struct {
  Id string `json:"id"`
  Name string `json:"name"`
  Devices []string `json:"devices"`
}

func Concat(a string, b string) string {
	var buffer bytes.Buffer
	buffer.WriteString(a)
	buffer.WriteString(b)
    return buffer.String()
}

func ListModels() []Model {
  data := make([]byte, 100)
  url := Concat(EI_URL, "model")
  resp, err := http.Get(url)
  if err == nil {
  	resp.Body.Read(data)
    logging.Debug(string(data))
    return nil
  } else {
  	logging.Debug("An error happened, %s", err.Error())
  	return nil
  }
}

/*
func SearchModels(creds) []Model {

}
*/