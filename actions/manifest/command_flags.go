package manifest

import "github.com/urfave/cli"

/*
	Current command flags list
*/

func commandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "template",
			EnvVar:      "TEMPLATE_PATH",
			Usage:       "filename / path for template generation",
			Required:    true,
		},
	}
}
