package manifests

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"pancake"
	"pancake/internal"
	"path"
	"strings"
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
			return fmt.Errorf("create temp direcotry: %s", err)
		}
		defer os.RemoveAll(tempDir)

		if err := internal.PullManifestTemplates(tempDir); err != nil {
			return fmt.Errorf("pull k8s manifests templates from '%s': %s", tempDir, err)
		}
		manifestsDir = tempDir
	} else {
		manifestsDir = ctx.String(flagPath)
	}

	// -- >
	var template pancake.K8STemplate
	raw, err := internal.ReadYaml(ctx.String(flagConfigs), template)
	if err != nil {
		return fmt.Errorf("yaml configs '%s': %s", ctx.String(flagConfigs), err)
	}

	fmt.Println("raw", raw, ctx.String(flagConfigs))
	fmt.Println()
	if err := pancake.Validate(raw.(pancake.K8STemplate)); err != nil {
		return err
	}

	manifests, err := pancake.GenerateManifest(raw.(pancake.K8STemplate), path.Join(manifestsDir, "k8s-templates"))
	if err != nil {
		return fmt.Errorf("generate manifests: %s", err)
	}
	// -- >

	if ctx.Bool(flagStdOut) {
		fmt.Println(strings.Join(manifests, "\n"))
	} else {
		fmt.Println("TODO: write to files")
	}

	log.Println("manifest generated")
	return nil
}
