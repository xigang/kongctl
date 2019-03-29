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
)

//docs: https://docs.konghq.com/1.0.x/admin-api/#route-object

// The Route entities defines rules to match client requests. Each Route is associated with a Service,
// and a Service may have multiple Routes associated to it.
// Every request matching a given Route will be proxied to its associated Service.

const (
	ROUTE_RESOURCE_OBJECT = "routes"
)

type RouteList struct {
	Data []RouteConfig `json:"data"`
}

type RouteConfig struct {
	ID            string    `json:"id"`                       //The route id
	Protocols     []string  `json:"protocols"`                //A list of the protocols this Route should allow
	Methods       []string  `json:"methods,omitempty"`        //A list of HTTP methods that match this route
	Hosts         []string  `json:"hosts,omitempty"`          //A list of domain names that match this route. When using http or https protocols, at least one of hosts, paths, or methods must be set
	Paths         []string  `json:"paths,omitempty"`          //A list of paths that match this route.When using http or https protocols, at least one of hosts, paths, or methods must be set.
	RegexPriority int       `json:"regex_priority,omitempty"` //A number used to choose which route resolves a given request when several routes match it using regexes simultaneously.
	StripPath     bool      `json:"strip_path,omitempty"`     //When matching a Route via one of the paths, strip the matching prefix from the upstream request URL
	PreserveHost  bool      `json:"preserve_host,omitempty"`  //When matching a Route via one of the hosts domain names, use the request Host header in the upstream request headers
	Service       ServiceID `json:"service"`                  //The Service this Route is associated to.
}

type ServiceID struct {
	ID string `json:"id"`
}

var routeCommonFlags = []cli.Flag{
	cli.StringSliceFlag{
		Name:  "protocols",
		Usage: "A list of the protocols this route should allow",
	},
	cli.StringSliceFlag{
		Name:  "methods",
		Usage: "A list of HTTP methods that match this Route",
	},
	cli.StringSliceFlag{
		Name:  "hosts",
		Usage: "A list of domain names that match this route",
	},
	cli.StringSliceFlag{
		Name:  "paths",
		Usage: "A list of paths that match this route",
	},
	cli.IntFlag{
		Name:  "regex_priority",
		Value: 0,
		Usage: "Determines the relative order of this Route against others when evaluating regex paths",
	},
	cli.BoolFlag{
		Name:  "strip_path",
		Usage: "When matching a route via one of the paths, strip the matching prefix from the upstream request URL",
	},
	cli.BoolFlag{
		Name:  "preserve_host",
		Usage: "When matching a route via one of the hosts domain names, use the request Host header in the upstream request headers",
	},
	cli.StringFlag{
		Name:  "service_id",
		Usage: "The service id this route is associated to",
	},
	cli.StringSliceFlag{
		Name:  "snis",
		Usage: "A list of SNIs that match this route when using stream routing",
	},
	cli.StringSliceFlag{
		Name:  "sources",
		Usage: "A list of IP sources of incoming connections that match this route when using stream routing",
	},
	cli.StringSliceFlag{
		Name:  "destinations",
		Usage: "A list of IP destinations of incoming connections that match this route when using stream routing",
	},
}

var RouteResourceObjectCommand = cli.Command{
	Name:  "route",
	Usage: "The kong route object.",
	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "create route object",
			Flags:  routeCommonFlags,
			Action: createRoute,
		},
		{
			Name:  "get",
			Usage: "retrieve route object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The route id",
				},
			},
			Action: getRoute,
		},
		{
			Name:  "delete",
			Usage: "delete route object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the route id",
				},
			},
			Action: deleteRoute,
		},
		{
			Name:  "list",
			Usage: "list all routes object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "size",
					Value: "100",
					Usage: "A limit on the number of objects to be returned per page",
				},
				cli.StringFlag{
					Name:  "offset",
					Usage: "A cursor used for pagination. offset is an object identifier that defines a place in the list",
				},
			},
			Action: getRoutes,
		},
	},
}

//createRoute create route
func createRoute(c *cli.Context) error {
	serviceID := c.String("service_id")
	if serviceID == "" {
		return fmt.Errorf("service_id is empty")
	}

	cfg := &RouteConfig{
		Protocols:     c.StringSlice("protocols"),
		Methods:       c.StringSlice("methods"),
		Hosts:         c.StringSlice("hosts"),
		Paths:         c.StringSlice("paths"),
		RegexPriority: c.Int("regex_priority"),
		StripPath:     c.Bool("strip_path"),
		PreserveHost:  c.Bool("preserve_host"),
		Service: ServiceID{
			ID: c.String("service_id"),
		},
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Post(ctx, ROUTE_RESOURCE_OBJECT, nil, cfg, nil)
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

//getRoute retrieve route object
func getRoute(c *cli.Context) error {
	id := c.String("id")

	var requestURL string
	if id != "" {
		requestURL = fmt.Sprintf("%s/%s", ROUTE_RESOURCE_OBJECT, id)
	} else {
		return fmt.Errorf("route id:%s is invlid", id)
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	q := url.Values{}
	q.Add("size", c.String("size"))

	if c.String("offset") != "" {
		q.Add("offset", c.String("offset"))
	}

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

func deleteRoute(c *cli.Context) error {
	id := c.String("id")

	var requestURL string
	if id != "" {
		requestURL = fmt.Sprintf("%s/%s", ROUTE_RESOURCE_OBJECT, id)
	} else {
		return fmt.Errorf("route id is empty")
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Delete(ctx, requestURL, nil, nil)
	if err != nil {
		return err
	}

	if serverResponse.StatusCode == http.StatusNoContent {
		fmt.Printf("delete route %s success\n", id)
	} else {
		return fmt.Errorf("failed to delete route %s.", id)
	}

	return nil
}

//getRoutes list all routes object.
func getRoutes(c *cli.Context) error {
	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Get(ctx, ROUTE_RESOURCE_OBJECT, nil, nil)
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
