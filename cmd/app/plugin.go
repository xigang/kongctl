package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/urfave/cli"

	"github.com/xigang/kongctl/common/client"
	"github.com/xigang/kongctl/common/tools"
	"github.com/xigang/kongctl/pkg/plugin/authentication"
	"github.com/xigang/kongctl/pkg/plugin/logging"
	"github.com/xigang/kongctl/pkg/plugin/utils"
)

// https://docs.konghq.com/0.14.x/admin-api/#plugin-object

// A Plugin entity represents a plugin configuration that will be executed during the HTTP request/response lifecycle.
// It is how you can add functionalities to Services that run behind Kong, like Authentication or Rate Limiting for example.
// You can find more information about how to install and what values each plugin takes by visiting the Kong Hub.

const (
	PLUGIN_RESOURCE_OBJECT = "plugins"
)

type Plugins struct {
	Data []CommonPluginConfig `json:"data"`
}

type CommonPluginConfig struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	RouteID    Route     `json:"route,omitempty"`
	ServiceID  ServiceID `json:"service,omitempty"`
	ConsumerID Comsumner `json:"consumer,omitempty"`
	Enabled    bool      `json:"enabled,omitempty"`
}

type Route struct {
	ID string `json:"id"`
}

type Comsumner struct {
	ID string `json:"id"`
}

var PluginResourceObjectCommand = cli.Command{
	Name:  "plugin",
	Usage: "The kong plugin object.",

	Subcommands: []cli.Command{
		{
			Name:   "avalible_plugins",
			Usage:  "current support plugins object",
			Action: avaiblePlugins,
		},
		{
			Name:  "create",
			Usage: "create a plugin object",
			Subcommands: []cli.Command{
				authentication.BasicAuthCommand,
				logging.StatsDCommand,
			},
		},
		{
			Name:  "get",
			Usage: "retrieve a plugin object",
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
			Usage: "list all plugins object",
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
			Usage: "delete a plugin object",
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
	fmt.Printf("%-20s\t%-20s\n", "name", "message")
	for name, info := range utils.AvaliblePlugins {
		fmt.Printf("%-20s\t%-20s\n", name, info)
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
	name := c.String("name")
	serviceID := c.String("service_id")
	routeID := c.String("route_id")
	consumerID := c.String("consumer_id")
	size := c.String("size")

	var requestURL string = fmt.Sprintf("%s", PLUGIN_RESOURCE_OBJECT)

	q := url.Values{}
	if name != "" {
		q.Add("name", name)
	}

	if serviceID != "" {
		q.Add("service_id", serviceID)
	}

	if routeID != "" {
		q.Add("route_id", routeID)
	}

	if consumerID != "" {
		q.Add("consumer_id", consumerID)
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

	plugins := &Plugins{}

	if err = json.Unmarshal(body, plugins); err != nil {
		return err
	}

	fmt.Printf("%-40s\t%-20s\t%-20s\n", "ID", "NAME", "ENABLED")
	for _, p := range plugins.Data {
		fmt.Printf("%-40s\t%-20s\t%-20t\n", p.ID, p.Name, p.Enabled)
	}

	// tools.IndentFromBody(body)
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
