package ei

import "testing"
import "github.com/nightrune/wrench/logging"

func TestRunDevice(t *testing.T) {
  logging.SetLoggingLevel(logging.LOG_DEBUG);
  logging.Info("Starting Test")
  ListModels();
  logging.Info("Ending Test")
}