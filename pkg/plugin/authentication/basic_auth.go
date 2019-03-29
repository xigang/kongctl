package authentication

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/urfave/cli"

	"github.com/xigang/kongctl/common/client"
	"github.com/xigang/kongctl/common/tools"
	"github.com/xigang/kongctl/pkg/plugin/utils"
)

//Basic Authentication
//https://docs.konghq.com/hub/kong-inc/basic-auth/

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

var BasicAuthCommand = cli.Command{
	Name: "basic-auth",
	Flags: append(utils.CommonPluginFlags, []cli.Flag{
		cli.BoolFlag{Name: "hide_credentials", Usage: "an optional boolean value telling the plugin to show or hide the credential from the upstream service"},
		cli.StringFlag{Name: "anonymous", Value: "", Usage: "an optional string (consumer uuid) value to use as an “anonymous” consumer if authentication fails"},
	}...),
	Usage: "create basic-auth plugin",
	Subcommands: []cli.Command{
		{
			Name:  "credential",
			Usage: "You can provision new username/password credentials by making the following HTTP request",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "username", Usage: "The username to use in the Basic Authentication，when consumer id is not empty"},
				cli.StringFlag{Name: "password", Usage: "The password to use in the Basic Authentication, when consumer id is not empty"},
			},
			Action: createBasicAuthCredential,
		},
	},
	Action: createBasicAuthPlugin,
}

func createBasicAuthPlugin(c *cli.Context) error {
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

type BasicAuthCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func createBasicAuthCredential(c *cli.Context) error {
	consumerID := c.String("consumer_id")
	username := c.String("username")
	password := c.String("password")

	if consumerID == "" || username == "" || password == "" {
		return fmt.Errorf("consumer: %s username: %s password: %s is not allow empty", consumerID, username, password)
	}

	cfg := &BasicAuthCredential{
		Username: username,
		Password: password,
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	requestURL := fmt.Sprintf("consumers/%s/basic-auth", consumerID)

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
