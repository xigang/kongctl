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

//docs: https://docs.konghq.com/0.14.x/admin-api/#target-object

// A target is an ip address/hostname with a port that identifies an instance of a backend service.
// Every upstream can have many targets, and the targets can be dynamically added. Changes are effectuated on the fly.

const (
	TARGET_RESOURCE_OBJECT = "targets"
)

type Target struct {
	Data []TargetConfig `json:"data"`
}

type TargetConfig struct {
	ID         string `json:"id,omitempty"`
	UpstreamID string `json:"upstream_id,omitempty"`
	// The target address (ip or hostname) and port. If omitted the port defaults to 8000. If the hostname resolves to an SRV record, the port value will overridden by the value from the dns record.
	Target string `json:"target"`
	//The weight this target gets within the upstream loadbalancer (0-1000, defaults to 100). If the hostname resolves to an SRV record, the weight value will overridden by the value from the dns record.
	Weight int `json:"weight,omitempty"`
}

var targetCommonFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "the target id",
	},
	cli.StringFlag{
		Name:  "upstream_id",
		Usage: "the upstream id",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "the upstream name",
	},
	cli.StringFlag{
		Name:  "target",
		Usage: "The target address (ip or hostname) and port,If omitted the port defaults to 8000",
	},
	cli.StringFlag{
		Name:  "weight",
		Value: "100",
		Usage: "The weight this target gets within the upstream loadbalancer (0-1000).",
	},
}

var TargetResourceObjectCommand = cli.Command{
	Name:  "target",
	Usage: "The kong target object.",

	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "Create target object",
			Flags:  targetCommonFlags,
			Action: createTarget,
		},
		{
			Name:   "list",
			Usage:  "Lists all targets currently active on the upstream’s load balancing wheel",
			Flags:  targetCommonFlags,
			Action: getTargets,
		},
		{
			Name:  "delete",
			Usage: "Disable a target in the load balancer",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the target id",
				},
				cli.StringFlag{
					Name:  "upstream_id",
					Usage: "the upstream id",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "the upstream name",
				},
				cli.StringFlag{
					Name:  "target",
					Usage: "The target address (ip or hostname) and port,If omitted the port defaults to 8000",
				},
			},
			Action: deleteTarget,
		},
	},
}

func createTarget(c *cli.Context) error {
	target := c.String("target")
	weight := c.Int("weight")
	id := c.String("upstream_id")
	name := c.String("name")

	if target == "" {
		return fmt.Errorf("the target is not allow empty")
	}

	var requestURL string
	if id != "" {
		requestURL = fmt.Sprintf("upstreams/%s/targets", id)
	} else if name != "" {
		requestURL = fmt.Sprintf("upstreams/%s/targets", name)
	} else {
		return fmt.Errorf("the upstream id and name is not empty")
	}

	cfg := &TargetConfig{
		Target: target,
		Weight: weight,
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

func getTargets(c *cli.Context) error {
	upstreamID := c.String("upstream_id")
	upstreamName := c.String("name")
	targetID := c.String("id")
	target := c.String("target")
	weight := c.String("weight")

	var requestURL string
	if upstreamID != "" {
		requestURL = fmt.Sprintf("upstreams/%s/targets", upstreamID)
	} else if upstreamName != "" {
		requestURL = fmt.Sprintf("upstreams/%s/targets", upstreamName)
	} else {
		return fmt.Errorf("the upstream id and name is not empty")
	}

	q := url.Values{}
	if targetID != "" {
		q.Add("id", targetID)
	}

	if target != "" {
		q.Add("target", target)
	}

	q.Add("weight", weight)

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

	targets := &Target{}
	if err = json.Unmarshal(body, targets); err != nil {
		return nil
	}

	fmt.Printf("%-35s\t%-35s\t%-20s\t%-20s\n", "ID", "UPSTREAM_ID", "TARGET", "WEIGHT")
	for _, v := range targets.Data {
		fmt.Printf("%-35s\t%-35s\t%-20s\t%-20d\n", v.ID, v.UpstreamID, v.Target, v.Weight)
	}

	return nil
}

func deleteTarget(c *cli.Context) error {
	upstreamID := c.String("upstream_id")
	upstreamName := c.String("name")
	target := c.String("target")
	targetID := c.String("id")

	if upstreamID == "" && upstreamName == "" {
		return fmt.Errorf("the upstream name and id is not allow empty")
	}

	if targetID == "" && target == "" {
		return fmt.Errorf("the target name and id is not allow empty")
	}

	var requestURL string
	if upstreamID != "" && targetID != "" {
		requestURL = fmt.Sprintf("/upstreams/%s/argets/%s", upstreamID, targetID)
	} else if upstreamID != "" && target != "" {
		requestURL = fmt.Sprintf("/upstreams/%s/argets/%s", upstreamID, target)
	} else if upstreamName != "" && targetID != "" {
		requestURL = fmt.Sprintf("/upstreams/%s/argets/%s", upstreamName, targetID)
	} else if upstreamName != "" && target != "" {
		requestURL = fmt.Sprintf("/upstreams/%s/argets/%s", upstreamName, target)
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	sereverResponse, err := client.GatewayClient.Delete(ctx, requestURL, nil, nil)
	if err != nil {
		return err
	}

	if sereverResponse.StatusCode == http.StatusNoContent {
		fmt.Printf("delete target success.")
	} else {
		return fmt.Errorf("delete target failed: %v", err)
	}
	return nil
}
