#include "notherlib.nut"
#include "logging.nut"

function test() {
  log("Help")
  var x = DoSomething(__LINE__)
  var y = __LINE__
}
