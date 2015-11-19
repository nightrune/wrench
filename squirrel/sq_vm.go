package squirrel

/*
#cgo CXXFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo CPPFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo CFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo LDFLAGS: -Lsquirrel/SQUIRREL3/lib -lsquirrel -lsqstdlib -lstdc++ -static-libstdc++
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>

#include "squirrel.h"
#include "sqstdio.h"
#include "sqstdaux.h"

#ifdef _MSC_VER
#pragma comment (lib ,"squirrel.lib")
#pragma comment (lib ,"sqstdlib.lib")
#endif

// sq_helpers.go
extern void SquirrelLog(char*);

void wrench_log(HSQUIRRELVM v, const SQChar *s, ...) {
  char buf[256];
  va_list arglist;
  va_start(arglist, s);
  vsnprintf(buf, 256, s, arglist);
  SquirrelLog(buf);
  va_end(arglist);
}

int device_run_script(const char* file_name) {
  HSQUIRRELVM v;
  v = sq_open(1024); // creates a VM with initial stack size 1024

  sqstd_seterrorhandlers(v);

  sq_setprintfunc(v, wrench_log, NULL); //sets the print function
  printf("Inside device_run_script\n");
  sq_pushroottable(v); //push the root table(were the globals of the script will be stored)
  // also prints syntax errors if any
  if (SQ_SUCCEEDED(sqstd_dofile(v, _SC(file_name), SQFalse, SQTrue)) == SQFalse) {
    sqstd_printcallstack(v);
    return -1;
  }

  sq_pop(v,1); //pops the root table
  sq_close(v);
  return 0;
}
*/
import "C"
import "unsafe"
import "github.com/nightrune/wrench/logging"

func RunScript(script_file string) {
	file_string := C.CString(script_file)
	rval := C.device_run_script(file_string)
	C.free(unsafe.Pointer(file_string))
	if rval != 0 {
		logging.Fatal("Squirrel script failed to run")
	}
}

/**
 * This will diverge from RunScript eventually
 */
func DeviceRunScript(device_file string) {
	RunScript(device_file)
}

/**
 * This will diverge from RunScript eventually
 */
func AgentRunScript(agent_file string) {
	RunScript(agent_file)
}
