package app

import (
	"github.com/urfave/cli"
)

var CertificateResourceObjectCommand = cli.Command{
	Name:  "certificate",
	Usage: "The kong certificate object.",

	Subcommands: []cli.Command{},
}
