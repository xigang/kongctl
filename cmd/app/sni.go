package app

import (
	"github.com/urfave/cli"
)

//TODO
var SNIResourceObjectCommand = cli.Command{
	Name:  "snis",
	Usage: "The kong sni object.",

	Subcommands: []cli.Command{},
}
