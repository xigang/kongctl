package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli"

	kongapp "github.com/xigang/kongctl/cmd/app"
	"github.com/xigang/kongctl/common/client"
)

func main() {
	app := cli.NewApp()
	app.Name = "kongctl"
	app.Usage = "kong(0.14.0) api gateway command line tool."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "http://127.0.0.1:8001",
			Usage: "api gateway(kong) server address",
		},
		cli.StringFlag{
			Name:  "auth",
			Value: "",
			Usage: "basic authoritarian for api gateway",
		},
	}

	app.Before = func(c *cli.Context) error {
		host := c.GlobalString("host")
		token := c.GlobalString("auth")

		if host == "" || token == "" {
			return fmt.Errorf("gateway(kong) host: %s or auth: %s is invalid", host, token)
		}

		customHTTPHeaders := make(map[string]string)
		customHTTPHeaders["Authorization"] = fmt.Sprintf("Basic %s", token)

		var err error
		if kongapp.GatewayClient, err = client.NewHTTPClient(host, customHTTPHeaders); err != nil {
			return err
		}

		return nil
	}

	app.Commands = []cli.Command{
		kongapp.ServiceCommand,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("failed to run kong api gateway command line tool.")
	}
}
