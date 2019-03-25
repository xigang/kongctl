package app

import (
	"github.com/urfave/cli"
)

var RouteResourceObjectCommand = cli.Command{
	Name:  "route",
	Usage: "The kong route object.",
	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "create route object",
			Flags:  []cli.Flag{},
			Action: createRouteObject,
		},
		{
			Name:   "get",
			Usage:  "retrieve route object",
			Flags:  []cli.Flag{},
			Action: getRouteObject,
		},
		{
			Name:   "delete",
			Usage:  "delete route object",
			Flags:  []cli.Flag{},
			Action: deleteRouteObject,
		},
		{
			Name:   "list",
			Usage:  "list all routes object",
			Flags:  []cli.Flag{},
			Action: getRoutesObject,
		},
		{
			Name:   "update",
			Usage:  "update route object",
			Flags:  []cli.Flag{},
			Action: updateRouteObject,
		},
	},
}

func createRouteObject(c *cli.Context) error {
	return nil
}

func getRouteObject(c *cli.Context) error {
	return nil
}

func deleteRouteObject(c *cli.Context) error {
	return nil
}

func getRoutesObject(c *cli.Context) error {
	return nil
}

func updateRouteObject(c *cli.Context) error {
	return nil
}
