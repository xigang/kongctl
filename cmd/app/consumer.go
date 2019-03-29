package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/urfave/cli"

	"github.com/xigang/kongctl/common/client"
	"github.com/xigang/kongctl/common/tools"
)

// https://docs.konghq.com/1.0.x/admin-api/#consumer-object

// The Consumer object represents a consumer - or a user - of a Service.
// You can either rely on Kong as the primary datastore,
// or you can map the consumer list with your database to keep consistency between Kong and your existing primary datastore.

const (
	CONSUMER_RESOURCE_OBJECT = "consumers"
)

type Consumer struct {
	Data []ConsumerConfig `json:"data"`
}

type ConsumerConfig struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	CustomID string `json:"custom_id"`
}

var ConsumerResourceObjectCommnad = cli.Command{
	Name:  "consumer",
	Usage: "The kong consumer object.",

	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create consumer object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username",
					Usage: "the unique username of the consumer",
				},
				cli.StringFlag{
					Name:  "custom_id",
					Usage: "the custom_id field for storing an existing unique ID for the consumer",
				},
			},
			Action: createConsumer,
		},
		{
			Name:   "list",
			Usage:  "list all consumers object",
			Flags:  []cli.Flag{},
			Action: getConsumers,
		},
		{
			Name:  "get",
			Usage: "retrieve consumer object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the consumer id",
				},
				cli.StringFlag{
					Name:  "username",
					Usage: "the consumer username",
				},
			},
			Action: getConsumber,
		},
		{
			Name:  "delete",
			Usage: "delete consumer object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "the consumer id",
				},
				cli.StringFlag{
					Name:  "username",
					Usage: "the consumer username",
				},
			},
			Action: deleteConsumber,
		},
	},
}

//createConsumer create a consumer resource object
func createConsumer(c *cli.Context) error {
	username := c.String("username")
	customID := c.String("custom_id")

	if username == "" && customID == "" {
		return fmt.Errorf("username: %s or custom id: %s invalid", username, customID)
	}

	cfg := &ConsumerConfig{
		Username: username,
		CustomID: customID,
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Post(ctx, CONSUMER_RESOURCE_OBJECT, nil, cfg, nil)
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

//getConsumers list all consumers resource object
func getConsumers(c *cli.Context) error {
	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Get(ctx, CONSUMER_RESOURCE_OBJECT, nil, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(serverResponse.Body)
	if err != nil {
		return err
	}

	var cm Consumer

	if err = json.Unmarshal(body, &cm); err != nil {
		return err
	}

	fmt.Printf("%-35s\t%-10s\t%-10s\n", "ID", "USERNAME", "CUSTOM_ID")
	for _, c := range cm.Data {
		fmt.Printf("%-35s\t%-10s\t%-10s\n", c.ID, c.Username, c.CustomID)
	}

	return nil
}

//getConsumber get a consumer resource object
func getConsumber(c *cli.Context) error {
	id := c.String("id")
	username := c.String("username")

	var requestURL string
	if id != "" {
		requestURL = fmt.Sprintf("%s/%s", CONSUMER_RESOURCE_OBJECT, id)
	} else if username != "" {
		requestURL = fmt.Sprintf("%s/%s", CONSUMER_RESOURCE_OBJECT, username)
	} else {
		return fmt.Errorf("username and id invalid.")
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

//deleteConsumber delete a consumer resource object
func deleteConsumber(c *cli.Context) error {
	id := c.String("id")
	username := c.String("username")

	var requestURL string
	if id != "" {
		requestURL = fmt.Sprintf("%s/%s", CONSUMER_RESOURCE_OBJECT, id)
	} else if username != "" {
		requestURL = fmt.Sprintf("%s/%s", CONSUMER_RESOURCE_OBJECT, username)
	} else {
		return fmt.Errorf("username and id invalid.")
	}

	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := client.GatewayClient.Delete(ctx, requestURL, nil, nil)
	if err != nil {
		return err
	}

	if serverResponse.StatusCode == http.StatusNoContent {
		fmt.Printf("delete consumer success.")
	} else {
		return fmt.Errorf("failed to delete consumer success: %v", err)
	}

	return nil
}
