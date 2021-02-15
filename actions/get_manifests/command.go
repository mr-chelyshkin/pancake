package get_manifests

import (
	"github.com/urfave/cli"
	"log"
	"pancake/internal"
)

/*
get-manifests command.

	download k8s manifests templates to local.
*/

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "get-manifests",
		Usage:  "download k8s manifests templates to local",

		Flags:  append(flags, commandFlags()...),
		Action: func(ctx *cli.Context) error {return run(ctx)},
	}
}

// --- >
func run(ctx *cli.Context) error {
	if err := internal.PullManifestTemplates(ctx.String(flagPath)); err != nil {
		return err
	}

	log.Println("templates saved in ", ctx.String(flagPath))
	return nil
}
