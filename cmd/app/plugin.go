package app

import (
	"github.com/urfave/cli"
)

const (
	PLUGIN_RESOURCE_OBJECT = "plugins"
)

type PluginConfig struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Route    RouteID     `json:"route,omitempty"`
	Service  ServiceID   `json:"service,omitempty"`
	Consumer ConsumberID `json:"consumer,omitempty"`
	RunOn    string      `json:"run_on"`
	Enabled  bool        `json:"enabled,omitempty"`
	//Config options
}

type ServiceID struct {
	ID string `json:"id,omitempty"`
}

type RouteID struct {
	ID string `json:"id,omitempty"`
}

type ConsumberID struct {
	ID string `json:"id,omitempty"`
}

var pluginCommonPlugin = []cli.Flag{
	cli.StringFlag{
		Name:  "name",
		Usage: "the plugin name",
	},
	cli.StringFlag{
		Name:  "id",
		Usage: "the plugin id",
	},
	cli.StringFlag{
		Name:  "route_id",
		Usage: "the unique identifier of the Route that should be associated to the newly-created plugin",
	},
	cli.StringFlag{
		Name:  "service_id",
		Usage: "the unique identifier of the Service that should be associated to the newly-created plugin",
	},
	cli.StringFlag{
		Name:  "consumer_id",
		Usage: "the unique identifier of the Consumer that should be associated to the newly-created plugin",
	},
	cli.StringFlag{
		Name:  "run_on",
		Value: "first",
		Usage: "control on which Kong nodes this plugin will run",
	},
	cli.BoolFlag{
		Name:  "enabled",
		Usage: "whether the plugin is applied",
	},
}

var PluginResourceObjectCommand = cli.Command{
	Name:  "plugin",
	Usage: "The kong plugin object.",

	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "create a global plugin",
			Flags:  pluginCommonPlugin,
			Action: createGlobalPlugin,
		},
	},
}

func createPlugin(c *cli.Context) error {
	return nil
}
