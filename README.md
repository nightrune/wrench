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

  On windows wrench uses mingw, On MACOSX the xcode build tools should be suffecient.
  Once Squirrel is build just use go build in the project directory.

  You'll also need a tool called GPP. The original site is here: 
  http://en.nothingisreal.com/wiki/GPP
  
  I keep a fork here: https://github.com/nightrune/gpp
  
  There is currently no support for other platforms. It should build just
  as easily on linux as it does windows. If there are problems please submit
  pull requests, and let the project know by submitting a bug.
  
  We've had one person that built on OSX and had no problems.

To Build use:
Make sure you are in the root wrench directory and use the following command.
go build --ldflags '-extldflags "-static"'

This makes sure we don't have to distribute .a, or .dll files around with wrench

Exact Build Steps on MAC OSX using ports.

udo port install gpp
sudo port install go

mkdir ~/Documents/go.home
cd ~/Documents/go.home
export GOPATH=~/Documents/go.home

go get github.com/nightrune/wrench

cd src/github.com/nightrune/wrench/squirrel/
tar xvfz ~/Downloads/squirrel_3_0_7_stable.tar.gz
cd SQUIRREL3
make

cd ../..
go build

ls -l wrench

