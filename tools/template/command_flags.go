package template

import "github.com/urfave/cli"

/*
	Current command flags list
*/

func commandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "local_flag",
			EnvVar:      "local_flag",
			Usage:       "local_flag",
			Value:       "local_flag",
		},
	}
}
