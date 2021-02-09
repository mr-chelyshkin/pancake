package main

import (
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"log"
	"os"
)

const appSlug = "mr-chelyshkin/pancake"

// --
func update() {
	v := semver.MustParse(Version)

	latest, err := selfupdate.UpdateSelf(v, appSlug)
	if err != nil {
		log.Println("Binary update failed:", err)
		return
	}
	if !latest.Version.Equals(v) {
		log.Println("cli app was updated, run it again", latest.Version)
		log.Println("Release note:\n", latest.ReleaseNotes)
		os.Exit(0)
	}
	return
}
