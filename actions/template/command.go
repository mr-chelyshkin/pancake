package template

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"pancake"
	"pancake/internal"
	"path/filepath"
)

/*
template command.

	generate template file for manual filling configs.
	that config file use cli-app for generate k8s manifests.
*/

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "template",
		Usage:  "generate yaml format template for manual filling [used by cli-app for manifiest generating]",

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

	fileDir, _ := filepath.Split(ctx.String("file"))
	if fileDir == "" {
		fileDir = "./"
	}
	if ok, err := internal.IsWritable(fileDir); !ok {
		return err
	}

	file, err := os.Create(ctx.String("file"))
	if file != nil || err != nil {
		return err
	}
	defer file.Close()

	if err := ioutil.WriteFile(file.Name(), templateBytes, 0644); err != nil {
		return err
	}
	log.Println(ctx.String("file"), " created")
	return nil
}
