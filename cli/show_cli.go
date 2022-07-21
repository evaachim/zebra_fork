package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/project-safari/zebra"
	"github.com/spf13/cobra"
)

/*
Build a zebra inventory command line client:

zebra show servers
zebra show users
zebra show registrations
zebra show networks
zebra show switches
zebra show racks
zebra show labs
zebra show data-centers
zebra command will use the private RSA key that the user created before registration and authenticate so there is no need to login for the CLI
all show commands will support label filters
*/

const version = "unknown"

// ErrNoCmd returns an error if no command flag was passed into the cli
var ErrNoCmd = errors.New("input the zebra command")

func main() {
	var zebraCmd = &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "show zebra",
		Short:        "show zebra resources",
		Long:         "commands that will fetch and show zebra resources and users and use the RSA key for authentication",
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			resource, _ := cmd.Flags().GetString("resource")

			if resource != "" {
				fmt.Println("Need method to fetch resource")
				getResourceWithFlag(resource)
			} else {

				return ErrNoCmd
			}
			return nil
		},
	}

	zebraCmd.Flags().String("zebra_cli", path.Join(
		func() string {
			s, _ := os.Getwd()

			return s
		}(), "zebra_cli_inventory"),
		"root directory of the  cli inventory",
	)

	zebraCmd.PersistentFlags().String("resource", "", "Name of resource to fetch, use singular")

	if err := zebraCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getResourceWithFlag(flg string) zebra.Resource {
	term := strings.ToLower(flg)

	resource := new(zebra.Resource)

	if term == "dc" || term == "data-center" {
		term = "datacenter"
	}

	fmt.Printf("Searched for resource: %v", term)

	switch term {
	case "switch":
		fmt.Println("fetch switches")
	case "server":
		fmt.Println("fetch servers")
	case "network":
		fmt.Println("fetch servers")
	case "user":
		fmt.Println("fetch users")
	case "registration":
		fmt.Println("fetch registrations")
	case "rack":
		fmt.Println("fetch racks")
	case "datacenter":
		fmt.Println("fetch data centers")
	}

	return *resource
}

func GetUserKey() string {
	fmt.Println("What is your private key. Please enter to authenticate: ")

	var key string

	fmt.Scanln(&key)

	return key
}
