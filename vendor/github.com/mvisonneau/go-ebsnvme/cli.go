package main

import (
	"fmt"
	"github.com/urfave/cli"

	"github.com/mvisonneau/go-ebsnvme/pkg/ebsnvme"
)

var version = "<devel>"

const (
	usage = "go-ebsnvme <block_device> [--volume-id|--device-name]"
)

// runCli : Generates cli configuration for the application
func runCli() (c *cli.App) {
	c = cli.NewApp()
	c.Name = "go-ebsnvme"
	c.Version = version
	c.Usage = "Fetch information about AWS EBS NVMe volumes"
	c.UsageText = usage
	c.EnableBashCompletion = true

	c.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "volume-id, i",
			Usage: "only print the EBS volume-id",
		},
		cli.BoolFlag{
			Name:  "device-name, n",
			Usage: "only print the name of the block device",
		},
	}

	c.Action = func(ctx *cli.Context) error {
		if len(ctx.Args()) != 1 ||
			(ctx.Bool("volume-id") && ctx.Bool("device-name")) {
			return cli.NewExitError("Usage: "+usage, 1)
		}

		d, err := ebsnvme.ScanDevice(ctx.Args().First())
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return cli.NewExitError("", 1)
		}

		if ctx.Bool("volume-id") {
			fmt.Println(d.VolumeID)
			return nil
		}

		if ctx.Bool("device-name") {
			fmt.Println(d.Name)
			return nil
		}

		fmt.Println(d.VolumeID)
		fmt.Println(d.Name)
		return nil
	}

	return
}
