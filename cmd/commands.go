package main

import (
	"github.com/urfave/cli"
	"pancake/actions/download_manifests"
	"pancake/actions/generate_configs"
	"pancake/actions/generate_manifests"
)

/*
	List of cli commands
*/

func commands(flags []cli.Flag) []cli.Command {
	return []cli.Command{
		generate_configs.Init(flags),
		generate_manifests.Init(flags),
		download_manifests.Init(flags),
	}
}
