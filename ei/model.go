package ei

import (
	"net/http"
	"bytes"
	"github.com/nightrune/wrench/logging"
  "encoding/base64"
  "encoding/json"
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

type ModelList struct {
  Models []Model `json:"models"`
}

type BuildClient struct {
  creds string
  http_client *http.Client
}

func NewBuildClient(api_key string) *BuildClient {
  client := new(BuildClient)
  client.http_client = &http.Client{}
  cred_data := []byte(api_key)
  client.creds = base64.StdEncoding.EncodeToString(cred_data)
  return client
}

func Concat(a string, b string) string {
	var buffer bytes.Buffer
	buffer.WriteString(a)
	buffer.WriteString(b)
    return buffer.String()
}

func (m BuildClient) SetAuthHeader(request *http.Request) {
  request.Header.Set("Authorization", "Basic " + m.creds)
}

func (m *BuildClient) ListModels() []Model {
  data := make([]byte, 100)
  full_response := new(bytes.Buffer)
  url := Concat(EI_URL, "models")
  req, _ := http.NewRequest("GET", url, nil)
  m.SetAuthHeader(req)
  resp, err := m.http_client.Do(req)
  if err == nil {
    var n int
  	for err = nil; err == nil; n, err = resp.Body.Read(data) {
      full_response.Write(data[:n])
    }
    list := new(ModelList)
    if err := json.Unmarshal(full_response.Bytes(), list); err != nil {
      logging.Warn("Failed to unmarshal data from models.. %s", err.Error());
    }
    for _, m := range list.Models {
      logging.Debug("Name: %s", m.Name) 
    }
    
    return nil
  } else {
  	logging.Debug("An error happened during model get, %s", err.Error())
  	return nil
  }
}

/*
func SearchModels(creds) []Model {

}
*/