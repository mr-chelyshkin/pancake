package manifests

import "github.com/urfave/cli"

/*
	Current command flags list
*/

const (
	flagConfigs = "configs"
	flagStdOut  = "stdout"
	flagStage   = "stage"
	flagPath    = "path"
	flagFile    = "file"
)

func commandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        flagConfigs,
			EnvVar:      "MANIFESTS_CONFIGS",
			Usage:       "file with configs for generating k8s manifests",
			Required:    true,
		},
		cli.StringFlag{
			Name:        flagPath,
			EnvVar:      "MANIFESTS_PATH",
			Usage:       "path to manifests templates (default pull from cvs)",
			Required:    false,
		},
		cli.BoolFlag{
			Name:        flagStdOut,
			EnvVar:      "TEMPLATE_STDOUT",
			Usage:       "output configs to stdout",
		},
		cli.StringFlag{
			Name:        flagStage,
			EnvVar:      "TEMPLATE_STAGE",
			Usage:       "choose manifests stage",
			Required:    true,
		},
	}
}
