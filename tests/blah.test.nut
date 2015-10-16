#include "logging.nut"
#include "sq_unit.nut"

log("Blah test start")

console("tests\blah.test.nut:17")
local y = NewMockObject()
local f = y.NewMethod("dude2", 1)
y.ExpectMethodCall(f, [1])
y.ExpectMethodCall(f, [2])
y.dude2(1)
log("Blah test end")