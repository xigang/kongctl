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
)

//docs:https://docs.konghq.com/0.14.x/admin-api/#upstream-objects

// The upstream object represents a virtual hostname and can be used to loadbalance incoming requests over multiple services (targets).
// So for example an upstream named service.v1.xyz for a Service object whose host is service.v1.xyz.
// Requests for this Service would be proxied to the targets defined within the upstream.

const (
	UPSTREAM_RESOURCE_OBJECT = "upstreams"
)

type Upstream struct {
	Data []UpstreamConfig `json:"data"`
}

type UpstreamConfig struct {
	//The upstream ID
	ID string `json:"ID,omitempty"`
	//This is a hostname, which must be equal to the host of a Service.
	Name string `json:"name"`
	//The number of slots in the loadbalancer algorithm (10-65536, defaults to 1000).
	Slots int `json:"slots,omitempty"`
	//What to use as hashing input: none, consumer, ip, header, or cookie (defaults to none resulting in a weighted-round-robin scheme).
	HashOn string `json:"hash_on,omitempty"`
	// What to use as hashing input if the primary hash_on does not return a hash (eg. header is missing, or no consumer identified). One of: none, consumer, ip, header, or cookie (defaults to none, not available if hash_on is set to cookie).
	HashFallback string `json:"hash_fallback,omitempty"`
	//The header name to take the value from as hash input (only required when hash_on is set to header)
	HashOnHeader string `json:"hash_on_header,omitempty"`
	//The header name to take the value from as hash input (only required when hash_fallback is set to header).
	HashFallbackHeader string `json:"hash_fallback_header,omitempty"`
	//The cookie name to take the value from as hash input (only required when hash_on or hash_fallback is set to cookie). If the specified cookie is not in the request, Kong will generate a value and set the cookie in the response.
	HashOnCookie string `json:"hash_on_cookie,omitempty"`
	//The cookie path to set in the response headers (only required when hash_on or hash_fallback is set to cookie, defaults to "/")
	HashOnCookiePath string `json:"hash_on_cookie_path,omitempty"`
	//Target health check
	HealthChecks HealthChecks `json:"healthchecks,omitempty"`
}

type HealthChecks struct {
	Active  Active  `json:"active,omitempty"`
	Passive Passive `json:"passive,omitempty"`
}

type Active struct {
	//Socket timeout for active health checks (in seconds).
	Timeout int `json:"timeout,omitempty"`
	//Number of targets to check concurrently in active health checks.
	Concurrency int `json:"concurrency,omitempty"`
	//Path to use in GET HTTP request to run as a probe on active health checks.
	HTTPPath string `json:"http_path,omitempty"`
	//Health checks
	Healthy Healthy `json:"healthy,omitempty"`
	//Unhealthy checks
	Unhealthy Unhealthy `json:"unhealthy,omitempty"`
}

type Passive struct {
	//Health checks
	Healthy Healthy `json:"healthy,omitempty"`
	//Unhealthy checks
	Unhealthy Unhealthy `json:"unhealthy,omitempty"`
}

type Healthy struct {
	//Interval between active health checks for healthy targets (in seconds). A value of zero indicates that active probes for healthy targets should not be performed.
	Interval int `json:"interval,omitempty"`
	//An array of HTTP statuses to consider a success, indicating healthiness, when returned by a probe in active health checks.
	HTTPStatuses []int `json:"http_statuses,omitempty"`
	//Number of successes in active probes (as defined by healthchecks.active.healthy.http_statuses) to consider a target healthy.
	Successes int `json:"successes,omitempty"`
}

type Unhealthy struct {
	//Interval between active health checks for unhealthy targets (in seconds). A value of zero indicates that active probes for unhealthy targets should not be performed.
	Interval int `json:"interval,omitempty"`
	//An array of HTTP statuses to consider a failure, indicating unhealthiness, when returned by a probe in active health checks.
	HTTPStatuses []int `json:"http_statuses,omitempty"`
	//Number of TCP failures in active probes to consider a target unhealthy.
	TCPFailures int `json:"tcp_failures,omitempty"`
	//Number of timeouts in active probes to consider a target unhealthy.
	Timeouts int `json:"timeouts,omitempty"`
	//Number of HTTP failures in active probes (as defined by healthchecks.active.unhealthy.http_statuses) to consider a target unhealthy.
	HTTPFailures int `json:"http_failures"`
}

var upstreamCommonFlags = []cli.Flag{
	cli.StringFlag{Name: "name", Usage: "This is a hostname, which must be equal to the host of a Service."},
	cli.IntFlag{Name: "slots", Value: 1000, Usage: "The number of slots in the loadbalancer algorithm (10-65536)"},
	cli.StringFlag{Name: "hash_on", Value: "none", Usage: "What to use as hashing input: none, consumer, ip, header, or cookie"},
	cli.StringFlag{Name: "hash_fallback", Value: "none", Usage: "What to use as hashing input if the primary hash_on does not return a hash (eg. header is missing, or no consumer identified). One of: none, consumer, ip, header, or cookie"},
	cli.StringFlag{Name: "hash_on_header", Usage: "The header name to take the value from as hash input (only required when hash_on is set to header)."},
	cli.StringFlag{Name: "hash_fallback_header", Usage: "The header name to take the value from as hash input (only required when hash_fallback is set to header)."},
	cli.StringFlag{Name: "hash_on_cookie", Usage: "The cookie name to take the value from as hash input (only required when hash_on or hash_fallback is set to cookie). If the specified cookie is not in the request, Kong will generate a value and set the cookie in the response."},
	cli.StringFlag{Name: "hash_on_cookie_path", Value: "/", Usage: "The cookie path to set in the response headers (only required when hash_on or hash_fallback is set to cookie)"},
	cli.IntFlag{Name: "healthchecks_active_timout", Usage: "Socket timeout for active health checks (in seconds)."},
	cli.IntFlag{Name: "healthchecks_active_concurrency", Usage: "Number of targets to check concurrently in active health checks."},
	cli.StringFlag{Name: "healthchecks_active_http_path", Usage: "Path to use in GET HTTP request to run as a probe on active health checks."},
	cli.IntFlag{Name: "healthchecks_active_healthy_interval", Usage: "Interval between active health checks for healthy targets (in seconds). A value of zero indicates that active probes for healthy targets should not be performed."},
	cli.IntSliceFlag{Name: "healthchecks_active_healthy_http_statuses", Usage: "An array of HTTP statuses to consider a success, indicating healthiness, when returned by a probe in active health checks."},
	cli.IntFlag{Name: "healthchecks_active_unhealthy_interval", Usage: "Interval between active health checks for unhealthy targets (in seconds). "},
	cli.IntSliceFlag{Name: "healthchecks_active_unhealthy_http_statuses", Usage: "An array of HTTP statuses to consider a failure, indicating unhealthiness, when returned by a probe in active health checks."},
	cli.IntFlag{Name: "healthchecks_active_unhealthy_tcp_failures", Usage: "Number of TCP failures in active probes to consider a target unhealthy."},
	cli.IntFlag{Name: "healthchecks_active_unhealthy_timeouts", Usage: "Number of timeouts in active probes to consider a target unhealthy."},
	cli.IntSliceFlag{Name: "healthchecks_passive_healthy_http_statuses", Usage: "An array of HTTP statuses which represent healthiness when produced by proxied traffic, as observed by passive health checks."},
	cli.IntFlag{Name: "healthchecks_passive_healthy_successes", Usage: "Number of successes in proxied traffic (as defined by healthchecks.passive.healthy.http_statuses) to consider a target healthy, as observed by passive health checks."},
	cli.IntSliceFlag{Name: "healthchecks_passive_unhealthy_http_statuses", Usage: "An array of HTTP statuses which represent unhealthiness when produced by proxied traffic, as observed by passive health checks."},
	cli.IntFlag{Name: "healthchecks_passive_unhealthy_tcp_failures", Usage: "Number of TCP failures in proxied traffic to consider a target unhealthy, as observed by passive health checks."},
	cli.IntFlag{Name: "healthchecks_passive_unhealthy_timeouts", Usage: "Number of timeouts in proxied traffic to consider a target unhealthy, as observed by passive health checks."},
	cli.IntFlag{Name: "healthchecks_passive_unhealthy_http_failures", Usage: "Number of HTTP failures in proxied traffic (as defined by healthchecks.passive.unhealthy.http_statuses) to consider a target unhealthy, as observed by passive health checks."},
}

var UpstreamResourceObjectCommand = cli.Command{
	Name:  "upstream",
	Usage: "The kong upstream object.",

	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "create upstream object",
			Flags:  upstreamCommonFlags,
			Action: createUpstream,
		},
		{
			Name:  "get",
			Usage: "get upstream object",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "name", Usage: "the upstream name"},
				cli.StringFlag{Name: "id", Usage: "the upstream id"},
			},
			Action: getUpstream,
		},
		{
			Name:  "list",
			Usage: "list all upstream object",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "name", Usage: "the upstream name"},
				cli.StringFlag{Name: "id", Usage: "the upstream id"},
				cli.StringFlag{Name: "size", Value: "100", Usage: "A limit on the number of objects to be returned."},
			},
			Action: getUpstreams,
		},
		{
			Name:  "delete",
			Usage: "delete upstream object",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "name", Usage: "the upstream name"},
				cli.StringFlag{Name: "id", Usage: "the upstream id"},
			},
			Action: deleteUpstream,
		},
	},
}

func createUpstream(c *cli.Context) error {
	name := c.String("name")
	if name == "" {
		return fmt.Errorf("the upstream name is not allow empty")
	}

	cfg := UpstreamConfig{
		Name:               c.String("name"),
		Slots:              c.Int("slots"),
		HashOn:             c.String("hash_on"),
		HashFallback:       c.String("hash_fallback"),
		HashOnHeader:       c.String("hash_on_header"),
		HashFallbackHeader: c.String("hash_fallback_header"),
		HashOnCookie:       c.String("hash_on_cookie"),
		HashOnCookiePath:   c.String("hash_on_cookie_path"),
		HealthChecks: HealthChecks{
			Active: Active{
				Timeout:     c.Int("timeout"),
				Concurrency: c.Int("concurrency"),
				HTTPPath:    c.String("http_path"),
				Healthy: Healthy{
					Interval:     c.Int("interval"),
					HTTPStatuses: c.IntSlice("http_statuses"),
					Successes:    c.Int("successes"),
				},
				Unhealthy: Unhealthy{
					Interval:     c.Int("interval"),
					HTTPFailures: c.Int("http_failures"),
					HTTPStatuses: c.IntSlice("http_statuses"),
					Timeouts:     c.Int("timeouts"),
					TCPFailures:  c.Int("tcp_failures"),
				},
			},
			Passive: Passive{
				Healthy: Healthy{
					HTTPStatuses: c.IntSlice("http_statuses"),
					Successes:    c.Int("successes"),
				},
				Unhealthy: Unhealthy{
					HTTPStatuses: c.IntSlice("http_statuses"),
					TCPFailures:  c.Int("tcp_failures"),
					Timeouts:     c.Int("timeouts"),
					HTTPFailures: c.Int("http_failures"),
				},
			},
		},
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Post(ctx, UPSTREAM_RESOURCE_OBJECT, nil, cfg, nil)
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

func getUpstream(c *cli.Context) error {
	name := c.String("name")
	id := c.String("id")

	var requestURL string
	if name != "" {
		requestURL = fmt.Sprintf("%s/%s", UPSTREAM_RESOURCE_OBJECT, name)
	} else if id != "" {
		requestURL = fmt.Sprintf("%s/%s", UPSTREAM_RESOURCE_OBJECT, id)
	} else {
		return fmt.Errorf("the upstream name and id is not allow empty")
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

func getUpstreams(c *cli.Context) error {
	name := c.String("name")
	id := c.String("id")
	size := c.String("size")

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	q := url.Values{}

	if id != "" {
		q.Add("id", id)
	}

	if name != "" {
		q.Add("name", name)
	}

	q.Add("size", size)

	serverResponse, err := client.GatewayClient.Get(ctx, UPSTREAM_RESOURCE_OBJECT, q, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(serverResponse.Body)
	if err != nil {
		return err
	}

	upstreams := &Upstream{}
	if err = json.Unmarshal(body, upstreams); err != nil {
		return err
	}

	fmt.Printf("%-35s\t%-20s\t%-20s\t%-20s%-20s\t%-20s\n", "ID", "NAME", "HASH_ON", "HASH_FALLBACK", "HASH_ON_COOKIE_PATH", "SLOTS")
	for _, u := range upstreams.Data {
		fmt.Printf("%-35s\t%-20s\t%-20s\t%-20s%-20s\t%-20d\n", u.ID, u.Name, u.HashOn, u.HashFallback, u.HashOnCookiePath, u.Slots)
	}

	return nil
}

func deleteUpstream(c *cli.Context) error {
	name := c.String("name")
	id := c.String("id")

	var requestURL string
	if name != "" {
		requestURL = fmt.Sprintf("%s/%s", UPSTREAM_RESOURCE_OBJECT, name)
	} else if id != "" {
		requestURL = fmt.Sprintf("%s/%s", UPSTREAM_RESOURCE_OBJECT, id)
	} else {
		return fmt.Errorf("the upstream name and id is not allow empty")
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Delete(ctx, requestURL, nil, nil)
	if err != nil {
		return err
	}

	if serverResponse.StatusCode == http.StatusNoContent {
		fmt.Printf("delete upstream success")
	} else {
		return fmt.Errorf("delete upstream failed: %v", err)
	}
	return nil
}
