package app

import (
	"context"
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

type ServiceConfig struct {
	ID             string `json:"id"`              //the service id
	Name           string `json:name`              //the service name
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

var ServiceCommand = cli.Command{
	Name:  "service",
	Usage: "the kong service object.",
	Flags: []cli.Flag{
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
			Usage: "shorthand attribute to set protocol, host, port and path at once",
		},
	},
	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "create service object",
			Action: create,
		},
		{
			Name:   "list",
			Usage:  "list all services object",
			Action: list,
		},
		{
			Name:   "get",
			Usage:  "retrieve service object",
			Action: get,
		},
		{
			Name:   "update",
			Usage:  "update service object",
			Action: update,
		},
		{
			Name:   "delete",
			Usage:  "delete service object",
			Action: del,
		},
	},
}

func checkArgs(c *cli.Context) error {
	name := c.GlobalString("name")
	host := c.GlobalString("host")
	path := c.GlobalString("path")
	url := c.GlobalString("url")

	if name == "" || host == "" || path == "" || url == "" {
		return fmt.Errorf("name: %s host: %s path: %s url: %s is invalid.", name, host, path, url)
	}

	return nil
}

func create(c *cli.Context) error {
	err := checkArgs(c)
	if err != nil {
		return err
	}

	sConfig := &ServiceConfig{
		Name:           c.GlobalString("name"),
		Retries:        c.GlobalInt("retries"),
		Protocol:       c.GlobalString("procotol"),
		Host:           c.GlobalString("host"),
		Port:           c.GlobalInt("port"),
		Path:           c.GlobalString("path"),
		ConnectTimeout: c.GlobalInt("connect_timeout"),
		WriteTimeout:   c.GlobalInt("write_timeout"),
		ReadTimeout:    c.GlobalInt("read_timeout"),
		URL:            c.GlobalString("url"),
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := GatewayClient.PATCH(ctx, SERVICE_RESOURCE_OBJECT, nil, sConfig, nil)
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

	tools.IndentFromBody(body)
	return nil
}

func get(c *cli.Context) error {
	name := c.GlobalString("name")
	id := c.GlobalString("id")

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
	name := c.GlobalString("name")
	id := c.GlobalString("id")

	var requestURL string
	if name != "" {
		requestURL = fmt.Sprintf("%s/%s", SERVICE_RESOURCE_OBJECT, name)
	} else if id != "" {
		requestURL = fmt.Sprintf("%s/%s", SERVICE_RESOURCE_OBJECT, id)
	} else {
		return fmt.Errorf("name: %s id: %s is invalid", name, id)
	}

	sConfig := &ServiceConfig{
		Protocol:       c.GlobalString("protocol"),
		Host:           c.GlobalString("host"),
		Port:           c.GlobalInt("port"),
		Path:           c.GlobalString("path"),
		Retries:        c.GlobalInt("retries"),
		ConnectTimeout: c.GlobalInt("connect_timeout"),
		WriteTimeout:   c.GlobalInt("write_timeout"),
		ReadTimeout:    c.GlobalInt("read_timeout"),
		URL:            c.GlobalString("url"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	serverResponse, err := GatewayClient.PATCH(ctx, requestURL, nil, sConfig, nil)
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
	name := c.GlobalString("name")
	id := c.GlobalString("id")

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
