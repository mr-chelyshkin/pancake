package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"pancake/internal"
	"time"
)

/*
Application endpoint.
	Description:
		cli-app for generate K8S manifests from predefined config.
		in manifest generation process cli-app take templates from repo: @@.
	Options:
		--skip-update:  bool flag, if set: skipping update process
	Commands:
		template:       generate yaml format template for manual filling
		manifest:       generate k8s manifests from template
		get-manifests:  download manifests templates to local
*/

var Version   = "v0.0.1" // -ldflags change on build (set value from CI/CD tag)
var Usage     = "pancake"
var UsageText = "pancake [global_flags] [command] [command_flags]"

func main() {
	app := cli.NewApp()

	app.UsageText = UsageText
	app.Version   = Version
	app.Usage     = Usage

	app.Flags     = globalFlags()
	app.Commands  = commands(app.Flags)

	app.Before = func(ctx *cli.Context) error {
		if !ctx.Bool("skip-update") {
			timer := time.NewTimer(3 *time.Second)

			select {
			case <-timer.C:
				log.Println("update timeout, skip auto update")
				return nil
			case result := <-__update__():
				if result {
					os.Exit(0)
				}
				os.Exit(1)
			}

		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

// -- >
func __update__() <-chan bool {
	ch := make(chan bool)
	go func() {
		ch <-internal.Update(Version)
	}()
	return ch
}
