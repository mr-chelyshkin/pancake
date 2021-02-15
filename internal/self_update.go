package internal

import (
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"log"
	"pancake/globals"
)

/*
self-update process.

	by "rhysd/go-github-selfupdate/selfupdate" get from github app versions (releases) and check it with current.
	if versions not equals - update cli-app binary.

	cli-app github repository set in [ appSlug ] constant.
	cli-app versions set as pattern [ v[0-9]+.[0-9]+.[0-9] ]

	cli-app repo should start workflow for build binaries for different OS and arch
	ONLY on tags with pattern [ v[0-9]+.[0-9]+.[0-9] ].

	workflow in: .github/workflow/release.yaml
*/

// -- >
func Update(currentVersion string) bool {
	v := semver.MustParse(currentVersion[1:])

	latest, err := selfupdate.UpdateSelf(v, globals.AppSlug)
	if err != nil {
		log.Println("cli-app update failed: ", err)
		log.Println("you can use flag '--skip-update' for skipping")
		return false
	}
	if !latest.Version.Equals(v) {
		log.Println("cli-app was updated, new version: ", latest.Version)
		log.Println("release note:\n", latest.ReleaseNotes)
		log.Println("now, you can use it.")
		return true
	}
	return false
}
