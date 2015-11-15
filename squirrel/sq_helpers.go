package squirrel

/*
#cgo CXXFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo CPPFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo CFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo LDFLAGS: -LSQUIRREL3/lib -lsquirrel -lsqstdlib -lstdc++

*/
import "C"
import "github.com/nightrune/wrench/logging"

//export SquirrelLog
func SquirrelLog(s *C.char) {
  log_value := C.GoString(s)
  logging.Info(log_value)
}