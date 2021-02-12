package manifest

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"pancake"
	"pancake/internal"
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
	templateBytes, err := internal.ReadFile(ctx.String("template"))
	if err != nil {
		return err
	}

	var obj pancake.K8STemplate
	k8sTemplate := yaml.Unmarshal(*templateBytes, obj)

	fmt.Println(k8sTemplate)

	return nil
}
