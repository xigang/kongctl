package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/urfave/cli"

	"github.com/xigang/kongctl/common/client"
	"github.com/xigang/kongctl/common/tools"
)

//docs: https://docs.konghq.com/0.14.x/admin-api/#service-object

// Service entities, as the name implies, are abstractions of each of your own upstream services.
// Examples of Services would be a data transformation microservice, a billing API, etc.
// The main attribute of a Service is its URL (where Kong should proxy traffic to),
// which can be set as a single string or by specifying its protocol, host, port and path individually.
// Services are associated to Routes (a Service can have many Routes associated with it).
// Routes are entry-points in Kong and define rules to match client requests. Once a Route is matched, Kong proxies the request to its associated Service

const (
	SERVICE_RESOURCE_OBJECT = "services"
)

type ServiceList struct {
	Data []ServiceConfig `json:"data"`
}

type ServiceConfig struct {
	ID             string `json:"id"`              //the service id
	Name           string `json:"name"`            //the service name
	Retries        int    `json:"retries"`         //the number if retries to execute upon failure to proxy.
	Protocol       string `json:"protocol"`        //the protocol used to communicate with the upstream.
	Host           string `json:"host"`            //the host of the upstream server
	Port           int    `json:"port"`            //the upstream server port
	Path           string `json:"path"`            //the path to be used in requests to the upstream server
	ConnectTimeout int    `json:"connect_timeout"` //the timeout in millilsends for establishing a connection to the upstream server
	WriteTimeout   int    `json:"write_timeout"`   //the timeout in milliseconds between two successive write operations for transmitting a request to the upstream server.
	ReadTimeout    int    `json:"read_timeout"`    //the timeout in milliseconds between two successive read operations for transmitting a request to the upstream server
	URL            string `json:"url"`             //shorthand attribute to set protocol, host, port and path at once. This attribute is write-only
}

var serviceFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "the service id",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "the serevice name",
	},
	cli.IntFlag{
		Name:  "retries",
		Value: 5,
		Usage: "the number of retries to execute upon failure to proxy",
	},
	cli.StringFlag{
		Name:  "procotol",
		Value: "http",
		Usage: "the protocol used to communicate with the upstream",
	},
	cli.StringFlag{
		Name:  "host",
		Usage: "the host of the upstream server",
	},
	cli.IntFlag{
		Name:  "port",
		Value: 80,
		Usage: "the upstream server port",
	},
	cli.StringFlag{
		Name:  "path",
		Value: "",
		Usage: "the path to be used requests to the upstream",
	},
	cli.IntFlag{
		Name:  "connect_timeout",
		Value: 60000,
		Usage: "the timeout in milliseconds for establishing a connection to the upstream server",
	},
	cli.IntFlag{
		Name:  "write_timeout",
		Value: 60000,
		Usage: "the timeout in milliseconds between two successive write operations for transmitting a request to the upstream server",
	},
	cli.IntFlag{
		Name:  "read_timeout",
		Value: 60000,
		Usage: "the timeout in milliseconds between two successive read operations for transmitting a request to the upstream server",
	},
	cli.StringFlag{
		Name:  "url",
		Value: "",
		Usage: "shorthand attribute to set protocol, host, port and path at once",
	},
}

var ServiceResourceObjectCommand = cli.Command{
	Name:  "service",
	Usage: "The kong service object.",

	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "create service object",
			Flags:  serviceFlags,
			Action: createService,
		},
		{
			Name:  "get",
			Usage: "retrieve service object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the service id",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "the service name",
				},
			},
			Action: getService,
		},
		{
			Name:  "delete",
			Usage: "delete service object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the service id",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "the service name",
				},
			},
			Action: deleteService,
		},
		{
			Name:   "list",
			Usage:  "list all services object",
			Action: getAllServices,
		},
		{
			Name:  "routes",
			Usage: "list routes associated to a service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the service id",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "the service name",
				},
			},
			Action: getRoutesByService,
		},
	},
}

//createService create service.
func createService(c *cli.Context) error {
	cfg := &ServiceConfig{
		Name:           c.String("name"),
		Retries:        c.Int("retries"),
		Protocol:       c.String("procotol"),
		Host:           c.String("host"),
		Port:           c.Int("port"),
		Path:           c.String("path"),
		ConnectTimeout: c.Int("connect_timeout"),
		WriteTimeout:   c.Int("write_timeout"),
		ReadTimeout:    c.Int("read_timeout"),
		URL:            c.String("url"),
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Post(ctx, SERVICE_RESOURCE_OBJECT, nil, cfg, nil)
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

//getAllServices list all services
func getAllServices(c *cli.Context) error {
	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Get(ctx, SERVICE_RESOURCE_OBJECT, nil, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(serverResponse.Body)
	if err != nil {
		return err
	}

	var services ServiceList
	if err = json.Unmarshal(body, &services); err != nil {
		return err
	}

	fmt.Printf("%-35s\t%-10s\t%-10s\t%-10s\t%-10s\t%-10s\t%-10s\t%-10s\n", "ID", "NAME", "PROCOTOL", "HOST", "PORT", "PATH", "READ_TIMEOUT", "WRITE_TIMEOUT")
	for _, s := range services.Data {
		fmt.Printf("%-35s\t%-10s\t%-10s\t%-10s\t%-10d\t%-10s\t%-10d\t%-10d\n", s.ID, s.Name, s.Protocol, s.Host, s.Port, s.Path, s.ReadTimeout, s.WriteTimeout)
	}
	return nil
}

//getService retrieve a service by name or id.
func getService(c *cli.Context) error {
	name := c.String("name")
	id := c.String("id")

	var requestURL string
	if name != "" {
		requestURL = fmt.Sprintf("%s/%s", SERVICE_RESOURCE_OBJECT, name)
	} else if id != "" {
		requestURL = fmt.Sprintf("%s/%s", SERVICE_RESOURCE_OBJECT, id)
	} else {
		return fmt.Errorf("name: %s or id: %s is invalid", name, id)
	}

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

//deleteService delete a service.
func deleteService(c *cli.Context) error {
	name := c.String("name")
	id := c.String("id")

	var requestURL string
	if name != "" {
		requestURL = fmt.Sprintf("%s/%s", SERVICE_RESOURCE_OBJECT, name)
	} else if id != "" {
		requestURL = fmt.Sprintf("%s/%s", SERVICE_RESOURCE_OBJECT, id)
	} else {
		return fmt.Errorf("name: %s id: %s is invalid", name, id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := client.GatewayClient.Delete(ctx, requestURL, nil, nil)
	if err != nil {
		return err
	}

	fmt.Printf("delete service %s success.\n", id)
	return nil
}

func getRoutesByService(c *cli.Context) error {
	serviceID := c.String("id")
	serviceName := c.String("name")

	var requestURL string
	if serviceID != "" {
		requestURL = fmt.Sprintf("%s/%s/routes", SERVICE_RESOURCE_OBJECT, serviceID)
	} else if serviceName != "" {
		requestURL = fmt.Sprintf("%s/%s/routes", SERVICE_RESOURCE_OBJECT, serviceName)
	} else {
		return fmt.Errorf("service id and anme is empty.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
