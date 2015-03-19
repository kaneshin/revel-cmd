package main

import (
	"strconv"

	"github.com/revel/revel"
	"github.com/revel/revel/harness"
)

var cmdPre = &Command{
	UsageLine: "pre [import path] [run mode] [port]",
	Short:     "pre a Revel application",
	Long: `
Run the Revel web application named by the given import path.

For example, to run the chat room sample application:

    revel pre github.com/revel/revel/samples/chat dev

The run mode is used to select which set of app.conf configuration should
apply and may be used to determine logic in the application itself.

Run mode defaults to "dev".

You can set a port as an optional third parameter.  For example:

    revel pre github.com/revel/revel/samples/chat prod 8080`,
}

func init() {
	cmdPre.Run = preApp
}

func preApp(args []string) {
	if len(args) == 0 {
		errorf("No import path given.\nRun 'revel help run' for usage.\n")
	}

	// Determine the run mode.
	mode := "dev"
	if len(args) >= 2 {
		mode = args[1]
	}

	// Find and parse app.conf
	revel.Init(mode, args[0], "")
	revel.LoadMimeConfig()

	// Determine the override port, if any.
	port := revel.HttpPort
	if len(args) == 3 {
		var err error
		if port, err = strconv.Atoi(args[2]); err != nil {
			errorf("Failed to parse port as integer: %s", args[2])
		}
		revel.HttpPort = port
	}

	revel.INFO.Printf("Running %s (%s) in %s mode\n", revel.AppName, revel.ImportPath, mode)
	revel.TRACE.Println("Base path:", revel.BasePath)

	// Else, just build and run the app.
	if err := harness.PreBuild(); err != nil {
		errorf("Failed to build app: %s", err)
	}
}
