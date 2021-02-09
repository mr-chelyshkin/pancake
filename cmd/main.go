package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

/*
Application endpoint.
	Description:
		cli-app for generate K8S manifests from predefined config.
		in manifest generation process cli-app take templates from repo: @@.

	Commands and Options:

*/

var Version   = "1.0.0" // -ldflags change on build (set value from CI/CD tag)
var Usage     = "pancake"
var UsageText = "pancake [command] [global_flags / command_flags]"

func main() {
	app := cli.NewApp()

	app.UsageText = UsageText
	app.Version   = Version
	app.Usage     = Usage

	app.Flags     = globalFlags()
	app.Commands  = commands(app.Flags)
	update()

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
