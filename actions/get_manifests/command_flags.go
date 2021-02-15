package get_manifests

import "github.com/urfave/cli"

/*
	Current command flags list
*/

const (
	defaultPath  = "./"
)

const (
	flagPath = "path"
)

func commandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        flagPath,
			EnvVar:      "MANIFESTS_PATH",
			Usage:       "path for download manifests templates",
			Value:       defaultPath,
			Required:    false,
		},
	}
}

