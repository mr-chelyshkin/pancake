package manifest

import (
	"github.com/urfave/cli"
)

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "manifest",
		Usage:  "generate manifests from template and manifests modules",

		Flags:  append(flags, commandFlags()...),
		Action: func(ctx *cli.Context) error {return run(ctx)},
	}
}

func run(ctx *cli.Context) error {



	return nil
}
