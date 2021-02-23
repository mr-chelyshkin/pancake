package generate_manifests

import "github.com/urfave/cli"

/*
	Current command flags list
*/

const (
	flagTemplates = "templates"
	flagConfigs   = "configs"
	flagStdOut    = "stdout"
	flagStage     = "stage"
	flagPath      = "path"
)

func commandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        flagStage,
			EnvVar:      "MANIFESTS_STAGE",
			Usage:       "choose manifests stage",
			Required:    true,
		},
		cli.StringFlag{
			Name:        flagConfigs,
			EnvVar:      "MANIFESTS_CONFIGS",
			Usage:       "file with configs for generating k8s manifests",
			Required:    true,
		},
		cli.StringFlag{
			Name:        flagTemplates,
			EnvVar:      "MANIFESTS_TEMPLATES",
			Usage:       "path to manifests templates (default pull from cvs)",
			Required:    false,
		},
		cli.StringFlag{
			Name:        flagPath,
			EnvVar:      "MANIFESTS_PATH",
			Usage:       "path to write k8s manifests",
			Required:    false,
		},
		cli.BoolFlag{
			Name:        flagStdOut,
			EnvVar:      "MANIFESTS_STDOUT",
			Usage:       "output configs to stdout",
		},
	}
}
