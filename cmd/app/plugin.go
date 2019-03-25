package app

import (
	"github.com/urfave/cli"
)

var PluginResourceObjectCommand = cli.Command{
	Name:  "plugin",
	Usage: "The kong plugin object.",

	Subcommands: []cli.Command{},
}
