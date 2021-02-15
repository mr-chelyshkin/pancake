package manifests

import (
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"pancake"
	"pancake/internal"
	"path"
)

/*
manifest command.

	generate k8s manifest files from k8s_template.yaml (generate from "template" command).
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
	var manifestsDir string

	if ctx.String(flagPath) == "" {
		tempDir, err := ioutil.TempDir("/tmp", "_manifest")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)

		if err := internal.PullManifestTemplates(tempDir); err != nil {
			return err
		}
		manifestsDir = tempDir
	} else {
		manifestsDir = ctx.String(flagPath)
	}

	// -- >
	var template pancake.K8STemplate
	raw, err := internal.ReadYaml(ctx.String(flagConfigs), template)
	if err != nil {
		return err
	}
	if err := pancake.GenerateManifest(raw.(pancake.K8STemplate), path.Join(manifestsDir, "k8s-templates")); err != nil {
		return err
	}

	log.Println("manifest generated")
	return nil
}
