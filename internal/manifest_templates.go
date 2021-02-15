package internal

import (
	"os/exec"
)

/*
clone k8s templates manifests process.

	Clone templates from project "manifestGitHTTPS" to income directory.
*/

const manifestGitHTTPS = "https://github.com/mr-chelyshkin/k8s-templates.git"

func PullManifestTemplates(dir string) error {
	cmd := exec.Command("git", "clone", manifestGitHTTPS)
	cmd.Dir = dir

	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}