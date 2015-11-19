package squirrel

import "testing"
import "github.com/nightrune/wrench/logging"

func TestRunDevice(t *testing.T) {
	logging.SetLoggingLevel(logging.LOG_DEBUG)
	logging.Info("Starting Test")
	RunScript("test.nut")
	logging.Info("Ending Test")
}
