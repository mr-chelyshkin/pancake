package main

import "github.com/urfave/cli"

/*
	List of global cli flags
*/

func globalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "global_flag",
			EnvVar:      "global_flag",
			Usage:       "global_flag",
			Value:       "global_shit",
		},
	}
}
