package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/revel/revel"
	"github.com/revel/revel/harness"
)

var cmdDev = &Command{
	UsageLine: "dev [import path] [target path] [custom build flag]",
	Short:     "dev a Revel application (e.g. for deployment)",
	Long: `
Build the Revel web application named by the given import path.
This allows it to be deployed and run on a machine that lacks a Go installation.

WARNING: The target path will be completely deleted, if it already exists!

For example:

    revel dev github.com/revel/revel/samples/chat /tmp/chat true
`,
}

func init() {
	cmdDev.Run = devApp
}

func devApp(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "%s\n%s", cmdDev.UsageLine, cmdDev.Long)
		return
	}

	appImportPath, destPath := args[0], args[1]
	if !revel.Initialized {
		revel.Init("", appImportPath, "")
	}

	custom := true
	if len(args) >= 3 {
		if args[3] == "false" {
			custom = false
		}
	}

	if rmerr := os.RemoveAll(destPath); rmerr != nil {
		errorf("Abort: %s does not look like a build directory.", destPath)
	}
	os.MkdirAll(destPath, 0777)

	app, reverr := harness.CustomBuild(custom)
	panicOnError(reverr, "Failed to build")

	destBinaryPath := path.Join(destPath, filepath.Base(app.BinaryPath))
	mustCopyFile(destBinaryPath, app.BinaryPath)
	mustChmod(destBinaryPath, 0755)
}
