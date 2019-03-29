package logging

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

//StatsD
//https://docs.konghq.com/hub/kong-inc/statsd/

// Log metrics for a Service, Route to a StatsD server. It can also be used to log metrics on Collectd daemon by enabling its Statsd plugin.

const (
	PLUGIN_STATSD = "statsd"
)

type Statsd struct {
	ID string `json:"ID,omitempty"`
	//The name of the plugin to use, in this case statsd
	Name string `json:"name"`
	//Consumer_id is the id of the Consumer we want to associate with this plugin.
	ConsumerID string       `json:"consumer_id,omitempty"`
	ServiceID  string       `json:"service_id,omitempty"`
	Enabled    bool         `json:"enabled,omitempty"`
	Config     StatsDConfig `json:"config"`
}

type StatsDConfig struct {
	//The IP address or host name to send data to.
	Host string `json:"host"`
	//The port to send data to on the upstream server
	Port int `json:"port"`
	//List of Metrics to be logged. Available values are described under Metrics.docs:https://docs.konghq.com/hub/kong-inc/statsd/#metrics
	Metrics []string `json:"metrics,omitempty"`
	//String to be prefixed to each metric’s name.
	Prefix string `json:"prefix,omitempty"`
}

var StatsDCommand = cli.Command{
	Name:  "statsd",
	Usage: "log metrics for a service, route to a StatsD server",
	Flags: append(utils.CommonPluginFlags, []cli.Flag{
		cli.StringFlag{Name: "host", Value: "127.0.0.1", Usage: "The IP address or host name to send data to"},
		cli.IntFlag{Name: "port", Value: 8125, Usage: "The port to send data to on the upstream server"},
		cli.StringFlag{Name: "prefix", Value: "kong", Usage: "String to be prefixed to each metric’s name."},
	}...),
	Action: createStatsDPlugin,
}

func createStatsDPlugin(c *cli.Context) error {
	name := c.String("name")
	host := c.String("host")
	port := c.Int("port")
	prefix := c.String("prefix")
	//TODO metrics use default data.

	serverID := c.String("service_id")
	routeID := c.String("route_id")
	consumerID := c.String("consumer_id")

	var requestURL string
	if serverID != "" {
		//Enabling the plugin on a Service
		requestURL = fmt.Sprintf("services/%s/plugins", serverID)
	} else if routeID != "" {
		//Enabling the plugin on a Route
		requestURL = fmt.Sprintf("routes/%s/plugins", routeID)
	} else {
		//global statsd plugin
		requestURL = "plugins"
	}

	config := StatsDConfig{
		Host:   host,
		Port:   port,
		Prefix: prefix,
	}

	statsd := Statsd{
		Name:       name,
		Config:     config,
		ConsumerID: consumerID,
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Post(ctx, requestURL, nil, statsd, nil)
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
