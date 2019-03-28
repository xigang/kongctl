package authentication

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/urfave/cli"

	"github.com/xigang/kongctl/common/client"
	"github.com/xigang/kongctl/common/tools"
)

//Basic Authentication

//Add Basic Authentication to a Service or a Route with username and password protection.
//The plugin will check for valid credentials in the Proxy-Authorization and Authorization header (in this order).

const (
	PLUGIN_BASIC_AUTH = "basic-auth"
)

type BasicAuthPluginConfig struct {
	Name string `json:"name"`
	//An optional boolean value telling the plugin to show or hide the credential from the upstream service. If true, the plugin will strip the credential from the request (i.e. the Authorization header) before proxying it.
	HideCredentials bool `json:"hide_credentials"`
	//An optional string (consumer uuid) value to use as an “anonymous” consumer if authentication fails. If empty (default), the request will fail with an authentication failure 4xx. Please note that this value must refer to the Consumer id attribute which is internal to Kong, and not its custom_id.
	Anonymous string `json:"anonymous"`
}

func CreateBasicAuthPlugin(c *cli.Context) error {
	name := c.String("name")
	serviceID := c.String("service_id")
	routeID := c.String("route_id")

	if name == "" {
		return fmt.Errorf("plugin name not allow empty")
	}

	var requestURL string
	if serviceID != "" {
		//Enabling the plugin on a Service
		requestURL = fmt.Sprintf("services/%s/plugins", serviceID)
	} else if routeID != "" {
		//Enabling the plugin on a Route
		requestURL = fmt.Sprintf("routes/%s/plugins", routeID)
	} else {
		//globale plugin
		requestURL = "plugins"
	}

	cfg := BasicAuthPluginConfig{
		Name:            c.String("name"),
		HideCredentials: c.Bool("hide_credentials"),
		Anonymous:       c.String("anonymous"),
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Post(ctx, requestURL, nil, cfg, nil)
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
