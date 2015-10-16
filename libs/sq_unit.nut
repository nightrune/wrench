#ifndef _SQ_UNIT_NUT
#define _SQ_UNIT_NUT

#include "logging.nut"

function NewMockMethod(name, n_args, checkcall) {
  if (n_args == 1) {
    return function(arg1) {
      checkcall(name, [arg1])
    }
  } else if (n_args == 2) {
    return function(arg1, arg2) {}
  } else if (n_args == 3) {
    return function(arg1, arg2, arg3) {}
  } else if (n_args == 4) {
    return function(arg1, arg2, arg3, arg4) {}
  }
  return function() {}
}

function NewMockFunction(name, n_args) {
  if (n_args == 1) {
    return function(original_this, arg1) {}
  } else if (n_args == 2) {
    return function(original_this, arg1, arg2) {}
  } else if (n_args == 3) {
    return function(original_this, arg1, arg2, arg3) {}
  } else if (n_args == 4) {
    return function(original_this, arg1, arg2, arg3, arg4) {}
  }
  return function(original_this) {
  }
}

function NewMockObject() {
  return {
    name = "a",
    calls = {},
    
    function CheckMethodCall(method_name, args) {
      local call_list = null;
      console("Function: " + method_name + " called with args: " + args)
      try {
        call_list = this.calls[method_name]
      } catch (e) {
        log("Didn't find a list of expected calls...")
      }
      if (call_list == null) {
        log("No call list and got call...")
        assert(false)
      }
      
      console("Call list args " + call_list.len() + " args " + args.len())
      if (call_list.len() != args.len()) {
        log("Length of expected args and received args is different for function: " + method_name)
        assert(false)
      }
    }
    
    function NewMethod(name, n_args) {
      local f = NewMockMethod(name, n_args, this.CheckMethodCall)
      this[name] <- f.bindenv(this)
      return name
    }
    
    function ExpectMethodCall(name, args) {
      local call_list = null
      try {
        call_list = this.calls[name]
      } catch (e) {
        
      }
      if (call_list == null) {
        this.calls[name] <- []
        call_list = this.calls[name]
      }
      call_list.append(args)
      console("Added call to: " + name + " with args " + args)
    }
  }
}

#endif // _SQ_UNIT_NUT