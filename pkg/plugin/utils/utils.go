package utils

import (
	"github.com/urfave/cli"
)

var AvaliblePlugins map[string]string = map[string]string{
	"basic-auth": "The plugin will check for valid credentials in the Proxy-Authorization and Authorization header",
	"statsd":     "Log metrics for a Service, Route to a StatsD server",
}

var CommonPluginFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "the plugin id",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "the plugin name",
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
		Usage: "The unique identifier of the consumer that overrides the existing settings for this specific consumer on incoming requests",
	},
	cli.BoolFlag{
		Name:  "enabled",
		Usage: "whether the plugin is applied",
	},
}
