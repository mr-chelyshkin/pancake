package manifest

import (
	"fmt"
	"github.com/urfave/cli"
	"pancake"
	"pancake/internal"
)

/*
manifest command.

	generate k8s manifest files from template.
*/

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "manifest",
		Usage:  "generate manifests from template and manifests modules",

		Flags:  append(flags, commandFlags()...),
		Action: func(ctx *cli.Context) error {return run(ctx)},
	}
}

// --- >
func run(ctx *cli.Context) error {
	var template pancake.K8STemplate

	raw, err := internal.ReadYaml(ctx.String("template"), template)
	if err != nil {
		return err
	}
	template = raw.(pancake.K8STemplate)


	fmt.Println(raw.(pancake.K8STemplate))

	return nil
}
