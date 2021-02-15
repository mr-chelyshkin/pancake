package main

import (
	"github.com/urfave/cli"
	"pancake/actions/get_manifests"
	"pancake/actions/manifests"
	"pancake/actions/template"
)

/*
	List of cli commands
*/

func commands(flags []cli.Flag) []cli.Command {
	return []cli.Command{
		template.Init(flags),
		manifests.Init(flags),
		get_manifests.Init(flags),
	}
}
