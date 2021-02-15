package template

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"log"
	"pancake"
	"pancake/internal"
)

/*
template command.

	generate template file for manual filling configs.
	that config file use cli-app for generate k8s manifests.
*/

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "template",
		Usage:  "generate yaml format configs for manual filling [used by cli-app for manifiest generating]",

		Flags:  append(flags, commandFlags()...),
		Action: func(ctx *cli.Context) error { return run(ctx) },
	}
}

// --- >
func run(ctx *cli.Context) error {
	template := pancake.GenerateTemplateObject(ctx.Int("apps"))
	templateBytes, err := yaml.Marshal(&template)
	if err != nil {
		return err
	}

	if ctx.Bool("stdout") {
		fmt.Println(string(templateBytes))
	}

	if err := internal.WriteFile(ctx.String("file"), templateBytes); err != nil {
		return err
	}

	log.Println(ctx.String("file"), " created")
	return nil
}
