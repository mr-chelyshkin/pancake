package generate_manifests

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"pancake"
	"pancake/globals"
	"pancake/internal"
	"path"
)

/*
manifest command.

	generate k8s manifest files from k8s_template.yaml (generate from "template" command).
*/

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "generate-manifests",
		Usage:  "generate manifests from template and manifests modules",

		Flags:  append(flags, commandFlags()...),
		Action: func(ctx *cli.Context) error {return run(ctx)},
	}
}

// -- >
func run(ctx *cli.Context) error {
	var manifestsDir string

	// -- >
	if ctx.String(flagTemplates) == "" {
		// create temp directory for manifests templates
		tempDir, err := ioutil.TempDir("/tmp", "_manifest")
		if err != nil {
			return fmt.Errorf("create temp direcotry: %s", err)
		}
		defer os.RemoveAll(tempDir)

		// pull manifests templates to temp directory
		if err := internal.PullManifestTemplates(tempDir); err != nil {
			return fmt.Errorf("pull k8s manifests templates from '%s': %s", tempDir, err)
		}

		manifestsDir = path.Join(tempDir, globals.ManifestGitProject)
	} else {
		manifestsDir = ctx.String(flagTemplates)
	}
	// -- >

	var template pancake.K8STemplate
	raw, err := internal.ReadYaml(ctx.String(flagConfigs), template)
	if err != nil {
		return fmt.Errorf("yaml configs '%s': %s", ctx.String(flagConfigs), err)
	}

	// validate k8s templates
	if err := pancake.Validate(raw.(pancake.K8STemplate)); err != nil {
		return fmt.Errorf("config validation errors:\n%s", err)
	}

	// generate manifests
	manifests, err := pancake.GenerateManifests(raw.(pancake.K8STemplate), manifestsDir)
	if err != nil {
		return fmt.Errorf("generate manifests: %s", err)
	}

	// return results
	if ctx.Bool(flagStdOut) {
		for _, app := range *manifests {
			fmt.Println(app)
		}
	} else {
		for name, data := range *manifests {
			file := path.Join(ctx.String(flagPath), fmt.Sprintf("%s.yaml", name))
			if err := internal.WriteFile(file, []byte(data)); err != nil {
				return fmt.Errorf("write manifest to %s, %s", file, err)
			}
		}
	}

	return nil
}
