package app

import (
	"github.com/urfave/cli"
)

var SNIResourceObjectCommand = cli.Command{
	Name:  "target",
	Usage: "The kong sni object.",

	Subcommands: []cli.Command{},
}
