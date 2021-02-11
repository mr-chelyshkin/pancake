package manifest

import (
	"fmt"
	"github.com/urfave/cli"
)

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "smth",
		Usage:  "this is smth",

		Flags:  append(flags, commandFlags()...),
		Action: func(ctx *cli.Context) error {return run(ctx)},
	}
}

func run(ctx *cli.Context) error {
	fmt.Println("smth")
	return nil
}
