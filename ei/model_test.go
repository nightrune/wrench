package ei

import "testing"
import "github.com/nightrune/wrench/logging"

func TestRunDevice(t *testing.T) {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Test")
  client := NewBuildClient("1d2da2b7e4e35667283af41ba2458527")
  client.ListModels();
  logging.Info("Ending Test")
}