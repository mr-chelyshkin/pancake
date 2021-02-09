package template

import (
	"fmt"
	"github.com/urfave/cli"
)

func Init(flags []cli.Flag) cli.Command{
	return cli.Command{
		Name:   "other",
		Usage:  "this is other",

		Flags:  append(flags, commandFlags()...),
		Action: func(ctx *cli.Context) error { return run(ctx) },
	}
}


func run(ctx *cli.Context) error {

	fmt.Println(ctx.String("global_flag"))
	fmt.Println(ctx.String("local_flag"))


	fmt.Println("other")
	return nil
}

// -- >
