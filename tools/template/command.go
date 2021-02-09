package template

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
	"log"
)

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "other",
		Usage:  "this is other",

		Flags:  append(flags, commandFlags()...),
		Action: func(ctx *cli.Context) error { return run(ctx) },
	}
}

const version = "0.0.1"

func run(ctx *cli.Context) error {

	v := semver.MustParse(version)
	fmt.Println(v)

	latest, err := selfupdate.UpdateSelf(v, "mr-chelyshkin/pancake")
	if err != nil {
		log.Println("Binary update failed:", err)
		return nil
	}

	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		log.Println("Current binary is the latest version", version)
	} else {
		log.Println("Successfully updated to version", latest.Version)
		log.Println("Release note:\n", latest.ReleaseNotes)
	}


	//fmt.Println(ctx.String("global_flag"))
	//fmt.Println(ctx.String("local_flag"))
	//
	//
	//fmt.Println("other")
	return nil
}

// -- >
