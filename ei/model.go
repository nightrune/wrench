package ei

import (
	"net/http"
	"bytes"
	"github.com/nightrune/wrench/logging"
  "encoding/base64"
  "encoding/json"
  "errors"
  "net/http/httputil"
  "io/ioutil"
)

const EI_URL = "https://build.electricimp.com/v4/"
const MODELS_ENDPOINT = "models"
const MODELS_REVISIONS_ENDPOINT = "revisions"

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

type BuildError struct {
	Code string `json:"code"`
	ShortMessage string `json:"message_short"`
}

type CodeRevisionResponse struct {
  Success bool `json:"success"`
  Revisions CodeRevisionLong `json:"revision"`
  Error BuildError `json:"error"`
}

type CodeRevisionsResponse struct {
  Success bool `json:"success"`
  Revisions []CodeRevisionShort `json:"revisions"`
}

type CodeRevisionShort struct {
  Version int `json:"version"`
  CreatedAt string `json:"created_at"`
  ReleaseNotes string `json:"release_notes"`
}

type CodeRevisionLong struct {
  Version int `json:"version,omitempty"`
  CreatedAt string `json:"created_at,omitempty"`
  DeviceCode string `json:"device_code,omitempty"`
  AgentCode string `json:"agent_code,omitempty"`
  ReleaseNotes string `json:"release_notes,omitempty"`
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

func (m *BuildClient) _complete_request(method string,
	url string, data []byte)  ([]byte, error) {
  var req *http.Request
  if data != nil {
  	req, _ = http.NewRequest(method, url, bytes.NewBuffer(data))
  }  else {
  	req, _ = http.NewRequest(method, url, nil)
  }

  m.SetAuthHeader(req)
  req.Header.Set("Content-Type", "application/json")
  resp, err := m.http_client.Do(req)
  if err == nil {
    dump, err := httputil.DumpResponse(resp, true)
    logging.Debug(string(dump))
    full_response, err := ioutil.ReadAll(resp.Body)
    if err != nil {
    	return full_response, err
    }
    return full_response, nil
  } else {
  	return nil, err
  }	
}

func (m *BuildClient) ListModels() (*ModelList, error) {
  list := new(ModelList)
  full_resp, err := m._complete_request("GET", Concat(EI_URL, "models"), nil)
  if err != nil {
  	logging.Debug("An error happened during model get, %s", err.Error())
  	return list, err
  }
  
  if err := json.Unmarshal(full_resp, list); err != nil {
    logging.Warn("Failed to unmarshal data from models.. %s", err.Error());
    return list, err
  }

  return list, nil
}

func (m *BuildClient) GetCodeRevisionList(model_id string) (
	[]CodeRevisionShort, error) {
  var url bytes.Buffer
  resp := new(CodeRevisionsResponse)
  url.WriteString(EI_URL)
  url.WriteString(MODELS_ENDPOINT)
  url.WriteString("/")
  url.WriteString(model_id)
  url.WriteString("/")
  url.WriteString(MODELS_REVISIONS_ENDPOINT)
  full_resp, err := m._complete_request("GET", url.String(), nil)
  if err != nil {
  	logging.Debug("Failed to get code revisions: %s", err.Error())
  	return resp.Revisions, err
  }


  if err := json.Unmarshal(full_resp, resp); err != nil {
    logging.Warn("Failed to unmarshal data from code revision.. %s", err.Error());
    return resp.Revisions, err
  }
  
  if resp.Success == false {
  	return resp.Revisions, errors.New("Error When retriveing Code Revisions")
  }
  return resp.Revisions, nil
}


func (m *BuildClient) GetCodeRevision(model_id string, build_num string) (CodeRevisionLong, error) {
  var url bytes.Buffer
  resp := new(CodeRevisionResponse)
  url.WriteString(EI_URL)
  url.WriteString(MODELS_ENDPOINT)
  url.WriteString("/")
  url.WriteString(model_id)
  url.WriteString("/")
  url.WriteString(MODELS_REVISIONS_ENDPOINT)
  url.WriteString("/")
  url.WriteString(build_num)
  full_resp, err := m._complete_request("GET", url.String(), nil)
  if err != nil {
  	logging.Debug("Failed to get code revisions: %s", err.Error())
  	return resp.Revisions, err
  }

  if err := json.Unmarshal(full_resp, resp); err != nil {
    logging.Warn("Failed to unmarshal data from code revision.. %s", err.Error());
    return resp.Revisions, err
  }
  
  if resp.Success == false {
  	return resp.Revisions, errors.New("Error When retriveing Code Revisions")
  }
  return resp.Revisions, nil
}

func (m *BuildClient) UpdateCodeRevision(model_id string,
	request *CodeRevisionLong) (CodeRevisionLong, error) {
  var url bytes.Buffer
  resp := new(CodeRevisionResponse)
  url.WriteString(EI_URL)
  url.WriteString(MODELS_ENDPOINT)
  url.WriteString("/")
  url.WriteString(model_id)
  url.WriteString("/")
  url.WriteString(MODELS_REVISIONS_ENDPOINT)

  req_string, err := json.Marshal(request)
  logging.Debug("Request String for upload: %s", req_string)
  full_resp, err := m._complete_request("POST", url.String(), req_string)
  if err != nil {
  	logging.Debug("Failed to update code revisions: %s", err.Error())
  	return resp.Revisions, err
  }

  if err := json.Unmarshal(full_resp, resp); err != nil {
    logging.Warn("Failed to unmarshal data from code revision update.. %s", err.Error());
    return resp.Revisions, err
  }
  
  if resp.Success == false {
  	return resp.Revisions, errors.New("Error When retriveing Code Revisions")
  }
  return resp.Revisions, nil
}

/*
func SearchModels(creds) []Model {

}
*/