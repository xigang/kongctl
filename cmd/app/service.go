package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-toolkit/common"
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
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

var commonFlags = []cli.Flag{
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

var ServiceCommand = cli.Command{
	Name:  "service",
	Usage: "the kong service object.",

	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "create service object",
			Flags:  commonFlags,
			Action: create,
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
			Action: get,
		},
		{
			Name:   "update",
			Usage:  "update service object",
			Flags:  commonFlags,
			Action: update,
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
			Action: del,
		},
		{
			Name:   "list",
			Usage:  "list all services object",
			Action: list,
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

func create(c *cli.Context) error {
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

func list(c *cli.Context) error {
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

func get(c *cli.Context) error {
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

func update(c *cli.Context) error {
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

	fmt.Printf("id: %s name: %s, host: %s, port: %d procotol: %s, path: %s\n", id, name, c.String("host"), c.Int("port"), c.String("procotol"), c.String("path"))

	cfg := &ServiceConfig{
		Protocol:       c.String("protocol"),
		Host:           c.String("host"),
		Port:           c.Int("port"),
		Path:           c.String("path"),
		Retries:        c.Int("retries"),
		ConnectTimeout: c.Int("connect_timeout"),
		WriteTimeout:   c.Int("write_timeout"),
		ReadTimeout:    c.Int("read_timeout"),
		URL:            c.String("url"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	serverResponse, err := GatewayClient.PATCH(ctx, requestURL, nil, cfg, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(serverResponse.Body)
	if err != nil {
		return err
	}
	common.IndentFromBody(body)
	return nil
}

func del(c *cli.Context) error {
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

	fmt.Printf("delete service success.\n")
	return nil
}
