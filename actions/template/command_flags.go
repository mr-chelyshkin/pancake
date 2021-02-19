package template

import "github.com/urfave/cli"

/*
	Current command flags list
*/

const (
	defaultFileName  = "k8s_template.yaml"
	defaultAppsCount = 1
)

const (
	flagApps   = "apps"
	flagFile   = "file"
	flagStdOut = "stdout"
)

func commandFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:        flagApps,
			EnvVar:      "CONFIGS_APPS_COUNT",
			Usage:       "number of apps in generating configs",
			Value:       defaultAppsCount,
		},
		cli.StringFlag{
			Name:        flagFile,
			EnvVar:      "CONFIGS_FILE",
			Usage:       "configs filename",
			Value:       defaultFileName,
		},
		cli.BoolFlag{
			Name:        flagStdOut,
			EnvVar:      "CONFIGS_STDOUT",
			Usage:       "output configs to stdout",
		},
	}
}
