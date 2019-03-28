package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/urfave/cli"

	"github.com/xigang/kongctl/common/client"
	"github.com/xigang/kongctl/common/tools"
	"github.com/xigang/kongctl/pkg/plugin/authentication"
)

// https://docs.konghq.com/0.14.x/admin-api/#plugin-object

// A Plugin entity represents a plugin configuration that will be executed during the HTTP request/response lifecycle.
// It is how you can add functionalities to Services that run behind Kong, like Authentication or Rate Limiting for example.
// You can find more information about how to install and what values each plugin takes by visiting the Kong Hub.

const (
	PLUGIN_RESOURCE_OBJECT = "plugins"
)

var avaliblePlugins map[string]string = map[string]string{
	authentication.PLUGIN_BASIC_AUTH: "The plugin will check for valid credentials in the Proxy-Authorization and Authorization header",
}

type CommonPluginConfig struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	RouteID    Route     `json:"route,omitempty"`
	ServiceID  ServiceID `json:"service,omitempty"`
	ConsumerID Comsumner `json:"consume,omitempty"`
	Enabled    bool      `json:"enabled,omitempty"`
	RunOn      string    `json:"run_on"`
	Protocols  []string  `json:"protocols"`
	Tags       string    `json:"tags,omitempty"`
}

type Route struct {
	ID string `json:"id"`
}

type Comsumner struct {
	ID string `json:"id"`
}

var commonPluginFlags = []cli.Flag{
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
	cli.BoolFlag{
		Name:  "enabled",
		Usage: "whether the plugin is applied",
	},
	cli.StringFlag{
		Name:  "run_on",
		Value: "first",
		Usage: "control on which Kong nodes this plugin will run, given a Service Mesh scenario. Accepted values are: * first, meaning “run on the first Kong node that is encountered by the request”",
	},
	cli.StringFlag{
		Name:  "tags",
		Usage: "an optional set of strings associated with the Plugin, for grouping and filtering",
	},
}

var PluginResourceObjectCommand = cli.Command{
	Name:  "plugin",
	Usage: "The kong plugin object.",

	Subcommands: []cli.Command{
		{
			Name:   "avalible_plugins",
			Usage:  "list current support plugins",
			Action: avaiblePlugins,
		},
		{
			Name:  "create",
			Usage: "create a plugin object",
			Subcommands: []cli.Command{
				{
					Name: "basic-auth",
					Flags: append(commonPluginFlags, []cli.Flag{
						cli.BoolFlag{Name: "hide_credentials", Usage: "an optional boolean value telling the plugin to show or hide the credential from the upstream service"},
						cli.StringFlag{Name: "anonymous", Value: "", Usage: "an optional string (consumer uuid) value to use as an “anonymous” consumer if authentication fails"}}...),
					Usage:  "create basic-auth plugin",
					Action: authentication.CreateBasicAuthPlugin,
				},
			},
		},
		{
			Name:  "get",
			Usage: "retrieve a plugin",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the plugin id",
				},
			},
			Action: getPlugin,
		},
		{
			Name:  "list",
			Usage: "list all plugins",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "service_id",
					Usage: "the service object id",
				},
				cli.StringFlag{
					Name:  "route_id",
					Usage: "the route object id",
				},
				cli.StringFlag{
					Name:  "size",
					Value: "100",
					Usage: "limit on the number of objects to be returned",
				},
			},
			Action: getPlugins,
		},
		{
			Name:  "delete",
			Usage: "delete a plugin",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the plugin id",
				},
			},
			Action: deletePlugin,
		},
	},
}

//avaiblePlugins list all avaible plugins
func avaiblePlugins(c *cli.Context) error {
	fmt.Println("plugin_name:")
	for name, info := range avaliblePlugins {
		fmt.Printf("%-s\t%-s\n", name, info)
	}
	return nil
}

//getPlugin get a plugin
func getPlugin(c *cli.Context) error {
	id := c.String("id")

	if id == "" {
		return fmt.Errorf("plugin id is empty")
	}

	requestURL := fmt.Sprintf("%s/%s", PLUGIN_RESOURCE_OBJECT, id)

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Get(ctx, requestURL, nil, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(serverResponse.Body)
	if err != nil {
		return err
	}

	tools.IndentFromBody(body)

	return nil
}

//getPlugins get all plugin
func getPlugins(c *cli.Context) error {
	serviceID := c.String("service_id")
	routeID := c.String("route_id")
	size := c.String("size")

	var requestURL string = fmt.Sprintf("%s", PLUGIN_RESOURCE_OBJECT)

	q := url.Values{}
	if serviceID != "" {
		q.Add("service_id", serviceID)
	}

	if routeID != "" {
		q.Add("route_id", routeID)
	}

	if size != "" {
		q.Add("size", size)
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Get(ctx, requestURL, q, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(serverResponse.Body)
	if err != nil {
		return err
	}

	tools.IndentFromBody(body)
	return nil
}

//deletePlugin delete a plugin
func deletePlugin(c *cli.Context) error {
	id := c.String("id")
	if id == "" {
		return fmt.Errorf("plugin id is empty")
	}

	requestURL := fmt.Sprintf("%s/%s", PLUGIN_RESOURCE_OBJECT, id)

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Delete(ctx, requestURL, nil, nil)
	if err != nil {
		return err
	}

	if serverResponse.StatusCode == http.StatusNoContent {
		fmt.Printf("delete plugin %s success", id)
	} else {
		return fmt.Errorf("delete plugin %s failed.", id)
	}

	return nil
}
