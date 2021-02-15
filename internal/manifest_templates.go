package internal

import (
	"os/exec"
	"pancake/globals"
)

/*
clone k8s templates manifests process.

	Clone templates from project "manifestGitHTTPS" to income directory.
*/

func PullManifestTemplates(dir string) error {
	cmd := exec.Command("git", "clone", globals.ManifestGitHTTPS)
	cmd.Dir = dir

	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}