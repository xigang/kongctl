package app

import (
	"github.com/urfave/cli"
)

var ConsumerResourceObjectCommnad = cli.Command{
	Name:  "consumer",
	Usage: "The kong consumer object.",

	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "create consumer object",
			Flags:  []cli.Flag{},
			Action: createConsumerObject,
		},
		{
			Name:   "list",
			Usage:  "list all consumers object",
			Flags:  []cli.Flag{},
			Action: getConsumersObject,
		},
		{
			Name:   "get",
			Usage:  "retrieve consumer object",
			Flags:  []cli.Flag{},
			Action: getConsumberObject,
		},
		{
			Name:   "update",
			Usage:  "update consumer object",
			Flags:  []cli.Flag{},
			Action: updateConsumberObject,
		},
		{
			Name:   "delete",
			Usage:  "delete consumer object",
			Flags:  []cli.Flag{},
			Action: deleteConsumberObject,
		},
	},
}

func createConsumerObject(c *cli.Context) error {
	return nil
}

func getConsumersObject(c *cli.Context) error {
	return nil
}

func getConsumberObject(c *cli.Context) error {
	return nil
}

func updateConsumberObject(c *cli.Context) error {
	return nil
}

func deleteConsumberObject(c *cli.Context) error {
	return nil
}
