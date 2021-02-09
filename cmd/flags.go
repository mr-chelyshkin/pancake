package main

import "github.com/urfave/cli"

/*
	List of global cli flags
*/

func globalFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:   "skip-update",
			EnvVar: "SKIP_UPDATE",
			Usage:  "skip update cli-application",
		},
	}
}
