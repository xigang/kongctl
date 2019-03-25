package app

import (
	"github.com/urfave/cli"
)

var TargetResourceObjectCommand = cli.Command{
	Name:  "target",
	Usage: "The kong target object.",

	Subcommands: []cli.Command{},
}
