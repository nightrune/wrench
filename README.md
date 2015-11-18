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