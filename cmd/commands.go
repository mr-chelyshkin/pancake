package main

import (
	"github.com/urfave/cli"
	"pancake/tools/manifest"
	"pancake/tools/template"
)

/*
	List of cli commands
*/

func commands(flags []cli.Flag) []cli.Command {
	return []cli.Command{
		template.Init(flags),
		manifest.Init(flags),
	}
}
