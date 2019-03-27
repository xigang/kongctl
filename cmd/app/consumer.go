package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/urfave/cli"

	"github.com/xigang/kongctl/common/tools"
)

const (
	CONSUMER_RESOURCE_OBJECT = "consumers"
)

//https://docs.konghq.com/1.0.x/admin-api/#consumer-object
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
					Usage: "field for storing an existing unique ID for the consumer - useful for mapping Kong with users in your existing database",
				},
			},
			Action: createConsumerObject,
		},
		{
			Name:   "list",
			Usage:  "list all consumers object",
			Flags:  []cli.Flag{},
			Action: getConsumersObject,
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
			Action: getConsumberObject,
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
			Action: deleteConsumberObject,
		},
	},
}

func createConsumerObject(c *cli.Context) error {
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

	serverResponse, err := GatewayClient.Post(ctx, CONSUMER_RESOURCE_OBJECT, nil, cfg, nil)
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

func getConsumersObject(c *cli.Context) error {
	ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cannel()

	serverResponse, err := GatewayClient.Get(ctx, CONSUMER_RESOURCE_OBJECT, nil, nil)
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

func getConsumberObject(c *cli.Context) error {
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

func deleteConsumberObject(c *cli.Context) error {
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

	serverResponse, err := GatewayClient.Delete(ctx, requestURL, nil, nil)
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
