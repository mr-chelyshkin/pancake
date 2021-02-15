package template

import "github.com/urfave/cli"

/*
	Current command flags list
*/

const (
	defaultFileName  = "k8s_template.yaml"
	defaultAppsCount = 1
)

func commandFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:        "apps",
			EnvVar:      "CONFIGS_APPS_COUNT",
			Usage:       "number of apps in generating configs",
			Value:       defaultAppsCount,
		},
		cli.StringFlag{
			Name:        "file",
			EnvVar:      "CONFIGS_FILE",
			Usage:       "configs filename",
			Value:       defaultFileName,
		},
		cli.BoolFlag{
			Name:        "stdout",
			EnvVar:      "TEMPLATE_STDOUT",
			Usage:       "output configs to stdout",
		},
	}
}
