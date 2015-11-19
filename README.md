# wrench
Preprocessor and Build Tools for the Electric Imp Platform

Wrench has many sub commands

To build a file use
- wrench build

- This preprocesses the files and generates the agent.nut, and device.nut files

To run tests use
- wrench test

- This preprocesses all .test.nut files and runs them.


If you api key and model is set you can upload your generated code to a specific
model
- wrench upload

If you want to list the models for your api key
- wrench model list



For Development:
  To build wrench you'll need the squirrel source:
  The tarball can be found here:
  http://sourceforge.net/projects/squirrel/

  Its expected that the Squirrel code lives in <project dir>/squirrel/SQUIRREL3
  See squirrel/sq_vm.go for more details.

  You'll need to build squirrel first.
  Follow its build instructions for your platform.

  On windows wrench uses mingw
  Once Squirrel is build just use go build in the project directory.

  There is currently no support for other platforms. It should build just
  as easily on linux as it does windows. If there are problems please submit
  pull requests, and let the project know by submitting a bug.

To Build use:
go build --ldflags '-extldflags "-static"'

This makes sure we don't have to distribute .a, or .dll files around with wrench