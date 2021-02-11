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
			EnvVar:      "TEMPLATE_APPS_COUNT",
			Usage:       "number of apps in generated template",
			Value:       defaultAppsCount,
		},
		cli.StringFlag{
			Name:        "file",
			EnvVar:      "TEMPLATE_PATH",
			Usage:       "filename / path for template generation",
			Value:       defaultFileName,
		},
		cli.BoolFlag{
			Name:        "stdout",
			Usage:       "generate tamplate to stdout",
			EnvVar:      "TEMPLATE_STDOUT",
		},
	}
}
