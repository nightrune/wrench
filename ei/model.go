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
const DEVICES_ENDPOINT = "devices"

type DeviceListResponse struct {
  Success bool `json:"success"`
  Devices []Device `json:"devices"`
}

type Device struct {
  Id string `json:"id,omitempty"`
  Name string `json:"name,omitempty"`
  ModelId string `json:"model_id,omitempty"`
  PowerState string `json:"powerstate,omitempty"`
  Rssi int `json:"rssi,omitempty"`
  AgentId string `json:"agent_id,omitempty"`
  AgentStatus string `json:"agent_status,omitempty"`
}

type Model struct {
  Id string `json:"id,omitempty"`
  Name string `json:"name"`
  Devices []string `json:"devices,omitempty"`
}

type ModelList struct {
  Models []Model `json:"models"`
}

type ModelResponse struct {
  Model Model `json:"model"`
  Success bool `json:"success"`
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

func (m *BuildClient) CreateModel(new_model *Model) (*Model, error) {
  var url bytes.Buffer
  resp := new(ModelResponse)
  url.WriteString(EI_URL)
  url.WriteString(MODELS_ENDPOINT)

  req_string, err := json.Marshal(new_model)
  logging.Debug("Request String for upload: %s", req_string)  
  full_resp, err := m._complete_request("POST", url.String(), req_string)
  if err != nil {
    logging.Debug("An error happened during model creation, %s", err.Error())
    return &resp.Model, err
  }
  
  if err := json.Unmarshal(full_resp, resp); err != nil {
    logging.Warn("Failed to unmarshal data from model response.. %s", err.Error());
    return &resp.Model, err
  }

  return &resp.Model, nil
}

func (m *BuildClient) DeleteModel(model_id string) (error) {
  var url bytes.Buffer
  resp := new(ModelResponse)
  url.WriteString(EI_URL)
  url.WriteString(MODELS_ENDPOINT)
  url.WriteString("/")
  url.WriteString(model_id)

  full_resp, err := m._complete_request("DELETE", url.String(), nil)
  if err != nil {
    logging.Debug("An error happened during model deletion, %s", err.Error())
    return err
  }
  
  if err := json.Unmarshal(full_resp, resp); err != nil {
    logging.Warn("Failed to unmarshal data from model response.. %s", err.Error());
    return err
  }

  if resp.Success == false {
    return errors.New("Error When retriveing Code Revisions")
  }

  return nil
}

const MODELS_DEVICE_RESTART_ENDPOINT = "restart"
func (m *BuildClient) RestartModelDevices(model_id string) (error) {
  var url bytes.Buffer
  resp := new(ModelResponse)
  url.WriteString(EI_URL)
  url.WriteString(MODELS_ENDPOINT)
  url.WriteString("/")
  url.WriteString(model_id)
  url.WriteString("/")
  url.WriteString(MODELS_DEVICE_RESTART_ENDPOINT)

  full_resp, err := m._complete_request("POST", url.String(), nil)
  if err != nil {
    logging.Debug("An error happened during model restart, %s", err.Error())
    return err
  }
  
  if err := json.Unmarshal(full_resp, resp); err != nil {
    logging.Warn("Failed to unmarshal data from model response.. %s", err.Error());
    return err
  }

  if resp.Success == false {
    return errors.New("Error When retriveing Code Revisions")
  }

  return nil
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

func (m *BuildClient) GetDeviceList() ([]Device, error) {
  var url bytes.Buffer
  resp := new(DeviceListResponse)
  url.WriteString(EI_URL)
  url.WriteString(DEVICES_ENDPOINT)

  full_resp, err := m._complete_request("GET", url.String(), nil)
  if err != nil {
  	logging.Debug("Failed to get device list: %s", err.Error())
  	return resp.Devices, err
  }

  if err := json.Unmarshal(full_resp, resp); err != nil {
    logging.Warn("Failed to unmarshal data from code revision update.. %s", err.Error());
    return resp.Devices, err
  }
  
  if resp.Success == false {
  	return resp.Devices, errors.New("Error When retriveing Code Revisions")
  }
  return resp.Devices, nil
}

const DEVICES_LOG_ENDPOINT = "logs"

type DeviceLogEntry struct {
  Timestamp string `json:"timestamp"`
  Type string `json:"type"`
  Message string `json:"message"`
}

type DeviceLogResponse struct {
	Logs []DeviceLogEntry `json:"logs"`
	PollUrl string `json:"poll_url"`
	Success bool `json:"success"`
}

func (m *BuildClient) GetDeviceLogs(device_id string) ([]DeviceLogEntry, error) {
  var url bytes.Buffer
  resp := new(DeviceLogResponse)
  url.WriteString(EI_URL)
  url.WriteString(DEVICES_ENDPOINT)
  url.WriteString("/")
  url.WriteString(device_id)
  url.WriteString("/")
  url.WriteString(DEVICES_LOG_ENDPOINT)
  full_resp, err := m._complete_request("GET", url.String(), nil)
  if err != nil {
  	logging.Debug("Failed to get device logs: %s", err.Error())
  	return resp.Logs, err
  }

  if err := json.Unmarshal(full_resp, resp); err != nil {
    logging.Warn("Failed to unmarshal data from device logs.. %s", err.Error());
    return resp.Logs, err
  }
  
  if resp.Success == false {
  	return resp.Logs, errors.New("Error When retriveing device logs")
  }
  return resp.Logs, nil
}

/*
func SearchModels(creds) []Model {

}
*/