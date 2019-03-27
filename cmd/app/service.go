package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/urfave/cli"

	"github.com/xigang/kongctl/common/tools"
)

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

var serviceCommonFlags = []cli.Flag{
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
			Flags:  serviceCommonFlags,
			Action: createServiceObject,
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
			Action: getServiceObject,
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
			Action: deleteServiceObject,
		},
		{
			Name:   "list",
			Usage:  "list all services object",
			Action: getAllServices,
		},
	},
}

func checkArgs(c *cli.Context) error {
	name := c.String("name")
	url := c.String("url")

	if name == "" || url == "" {
		return fmt.Errorf("name: %s url: %s is invalid.", name, url)
	}

	return nil
}

//createServiceObject create service
func createServiceObject(c *cli.Context) error {
	err := checkArgs(c)
	if err != nil {
		return err
	}

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

	serverResponse, err := GatewayClient.Post(ctx, SERVICE_RESOURCE_OBJECT, nil, cfg, nil)
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

	serverResponse, err := GatewayClient.Get(ctx, SERVICE_RESOURCE_OBJECT, nil, nil)
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

	fmt.Printf("ID\t\t\t\t\t\t HOST_NAME\t\t\t PORT\t\t\t NAME\n")
	for _, s := range services.Data {
		fmt.Printf("id:%s\t\t host:%s\t\t port:%d\t\t name: %s\n", s.ID, s.Host, s.Port, s.Name)
	}
	return nil
}

//getAllServices retrieve a service
func getServiceObject(c *cli.Context) error {
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

	serverResponse, err := GatewayClient.Get(ctx, requestURL, nil, nil)
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

//deleteServiceObject delete a service
func deleteServiceObject(c *cli.Context) error {
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

	_, err := GatewayClient.Delete(ctx, requestURL, nil, nil)
	if err != nil {
		return err
	}

	fmt.Printf("delete service %s success.\n", id)
	return nil
}
