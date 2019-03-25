package app

import (
	"github.com/urfave/cli"
)

var UpstreamResourceObjectCommand = cli.Command{
	Name:  "upstream",
	Usage: "The kong upstream object.",

	Subcommands: []cli.Command{},
}
