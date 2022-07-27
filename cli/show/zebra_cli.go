package show

import (
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
zebra command will use the private RSA key that the user created
before registration and authenticate so there is no need to login for the CLI
all show commands will support label filters
*/

// create inventory of commands for the zebra cli.
func NewZebra() *cobra.Command {
	// default zebra command to show resources.
	zebraCmd := &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:     "show resources",
		Short:   "command to show zebra resources",
		Version: version,
	}

	zebraCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "users",
		Short:        "show zebra userss",
		RunE:         ShowUsr,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	zebraCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "regs",
		Short:        "show zebra registrations",
		RunE:         ShowReg,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	return zebraCmd
}

// for network resources.
func NewNetCmd() *cobra.Command {
	netCmd := &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "vlan-pools",
		Short:        "show zebra vlan pools",
		RunE:         ShowVlan,
		SilenceUsage: true,
	}

	netCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "switches",
		Short:        "show zebra switches",
		RunE:         ShowSw,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	zebraCmd := NewZebra()

	zebraCmd.AddCommand(netCmd)

	return zebraCmd
}

// for server resources.
func NewSrvCmd() *cobra.Command {
	srvCmd := &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "servers",
		Short:        "show zebra servers",
		RunE:         ShowServ,
		SilenceUsage: true,
	}

	srvCmd.AddCommand((&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "esx",
		Short:        "show zebra esx-servers",
		RunE:         ShowESX,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	}))

	srvCmd.AddCommand((&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "vcenter",
		Short:        "show zebra vcenter",
		RunE:         ShowVC,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	}))

	zebraCmd := NewZebra()

	zebraCmd.AddCommand(srvCmd)

	return zebraCmd
}

// for dc resources.
func NewDCCmd() *cobra.Command {
	dcCmd := &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "dc",
		Short:        "show datacenter information",
		RunE:         ShowDC,
		SilenceUsage: true,
	}

	dcCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "labs",
		Short:        "show zebra labss",
		RunE:         ShowLab,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	dcCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "racks",
		Short:        "show zebra racks",
		RunE:         ShowRack,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	zebraCmd := NewZebra()

	zebraCmd.AddCommand(dcCmd)

	return zebraCmd
}
