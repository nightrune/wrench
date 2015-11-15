package squirrel

/*
#cgo CXXFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo CPPFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo CFLAGS: -ISQUIRREL3/include -ISQUIRREL3
#cgo LDFLAGS: -LSQUIRREL3/lib -lsquirrel -lsqstdlib -lstdc++
#include <stdarg.h>
#include <stdio.h>

#include "squirrel.h"
#include "sqstdio.h"
#include "sqstdaux.h"


#ifdef _MSC_VER
#pragma comment (lib ,"squirrel.lib")
#pragma comment (lib ,"sqstdlib.lib")
#endif

// sq_helpers.go
extern void SquirrelLog(char*);

// This example calls a function inside the squirrel script
void call_foo(HSQUIRRELVM v, int n,float f,const SQChar *s) {
  int top = sq_gettop(v); //saves the stack size before the call
  sq_pushroottable(v); //pushes the global table
  sq_pushstring(v,_SC("foo"),-1);
  if(SQ_SUCCEEDED(sq_get(v,-2))) { //gets the field 'foo' from the global table
    sq_pushroottable(v); //push the 'this' (in this case is the global table)
    sq_pushinteger(v,n); 
    sq_pushfloat(v,f);
    sq_pushstring(v,s,-1);
    sq_call(v,4,0,0); //calls the function 
  }
  sq_settop(v,top); //restores the original stack size
}

void wrench_log(HSQUIRRELVM v, const SQChar *s, ...) { 
  char buf[256];
  va_list arglist;
  va_start(arglist, s); 
  sprintf(buf, s, arglist);
  SquirrelLog(buf);
  va_end(arglist); 
}

int device_run_script() {
  HSQUIRRELVM v; 
  v = sq_open(1024); // creates a VM with initial stack size 1024 

  sqstd_seterrorhandlers(v);

  sq_setprintfunc(v, wrench_log, NULL); //sets the print function

  sq_pushroottable(v); //push the root table(were the globals of the script will be stored)
  // also prints syntax errors if any 
  if (SQ_SUCCEEDED(sqstd_dofile(v, _SC("test.nut"), 0, 1)) == SQFalse) {
    return -1;
  }


  sq_pop(v,1); //pops the root table
  sq_close(v); 


  return 0; 
} 
*/
import "C"

func DeviceRunScript(device_file string) {
    C.device_run_script()
}