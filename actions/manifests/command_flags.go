package manifests

import "github.com/urfave/cli"

/*
	Current command flags list
*/

func commandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "configs",
			EnvVar:      "MANIFESTS_CONFIGS",
			Usage:       "file with configs for generating k8s manifests",
			Required:    true,
		},
		cli.StringFlag{
			Name:        "path",
			EnvVar:      "MANIFESTS_PATH",
			Usage:       "path to manifests templates (default pull from cvs)",
			Required:    false,
		},
	}
}
