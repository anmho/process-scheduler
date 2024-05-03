
Precompiled executable platform binaries contained in ./bin/project1 (darwin-arm64 and linux-amd64)
Output is in ./out/output.txt

Build instructions
    1) Prerequisite: Install go 1.20
    2) Dependencies: go mod download
    3) Build: `make` to build for the current platform
    4) Preinstalled build commands for Linux amd64 platforms: `make linux-amd64` and `make darwin-arm64`
To run locally:
    1) Run: ./bin/project1 < ./rsrc/input.txt > ./out/<output_name>.txt`
    2) Output will be in ./out/<output_name>.txt


File descriptions
    An explanation of the folders and files within the project. Here’s a possible breakdown based on the image:
.
├── Makefile -- Preset commands for building
├── README.txt -- README.txt
├── bin -- Precompiled executables for a few platforms
│   ├── project1-darwin-arm64
│   └── project1-linux-amd64
├── cmd -- Main function. Responsible for parsing commands and displaying outputs
│    └── main.go
├── go.mod
├── go.sum
├── internal -- Internal components for the program
│    ├── manager -- Manager package
│    │   ├── manager.go -- Contains functions related to the manager.
│    │   ├── manager_test.go -- Contains tests for the manager
│    │   ├── readylist.go -- Contains the readylist struct and related functions
│    │   ├── readylist_test.go -- Contains tests for the readylist
│    │   ├── resource.go -- Contains functions related to the resource management
│    │   ├── scheduler.go -- Contains scheduler function
│    │   ├── timer.go -- Contains timeout function
│    │   └── wait.go -- Contains functions to manage resource waitlist
│    ├── process
│    │   └── pcb.go -- Contains PCB struct and functions for the PCB
│    └── resources
│        ├── resource.go -- Contains RCB struct and related functions
│        └── waitlist.go -- Contains resource waitlist structs and related functions
├── out
│    └── output.txt
└── rsrc
    ├── input.txt
    ├── sample-input-SQ24.txt
    ├── sample-output.txt
    ├── test1.txt
    ├── test2.txt
    ├── test3.txt
    ├── test4.txt
    ├── test5.txt
    └── test6.txt
