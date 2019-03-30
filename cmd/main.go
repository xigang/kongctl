package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/golang/glog"
	"github.com/urfave/cli"

	kongapp "github.com/xigang/kongctl/cmd/app"
	"github.com/xigang/kongctl/common/client"
)

func main() {
	app := cli.NewApp()
	app.Name = "kongctl"
	app.Usage = "kong(0.14.0) api gateway command line tool.\n\t https://docs.konghq.com/0.14.x/admin-api"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			EnvVar: "KONG_HOST",
			Usage:  "api gateway(kong) server address",
		},
		cli.StringFlag{
			Name:   "auth",
			EnvVar: "KONG_AUTH",
			Usage:  "basic authoritarian for api gateway",
		},
	}

	app.Before = func(c *cli.Context) error {
		host := c.GlobalString("host")
		token := c.GlobalString("auth")

		if host == "" {
			fmt.Printf("please specify the KONG_HOST and KONG_AUTH environment variables")
		}

		customHTTPHeaders := make(map[string]string)
		customHTTPHeaders["Authorization"] = fmt.Sprintf("Basic %s", token)

		var err error
		if client.GatewayClient, err = client.NewHTTPClient(host, customHTTPHeaders); err != nil {
			return err
		}

		return nil
	}

	app.Commands = []cli.Command{
		kongapp.ServiceResourceObjectCommand,
		kongapp.RouteResourceObjectCommand,
		kongapp.ConsumerResourceObjectCommnad,
		kongapp.CertificateResourceObjectCommand,
		kongapp.PluginResourceObjectCommand,
		kongapp.SNIResourceObjectCommand,
		kongapp.UpstreamResourceObjectCommand,
		kongapp.TargetResourceObjectCommand,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		glog.Errorf("%+v", err)
	}
}
