package show

import (
	"errors"
	"os"
	"path"
	"path/filepath"

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
zebra command will use the private RSA key that the user created before
registration and authenticate so there is no need to login for the CLI
all show commands will support label filters
*/

// ErrNoCmd returns an error if no command flag was passed into the cli.
var ErrNoCmd = errors.New("input the zebra command")

// ErrWrongCmd returns an error if the command that was input is wrong or does not exist.
var ErrWrongCmd = errors.New("not a zebra command")

// version for the zebra tool.
var version = "unknown"

// rootCmd is the default command for the zebra cli.
func New() *cobra.Command {
	name := filepath.Base(os.Args[0])

	rootCmd := &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:     name,
		Short:   "zebra resource reservation client",
		Version: version + "\n",
	}
	rootCmd.SetVersionTemplate(version + "\n")
	rootCmd.PersistentFlags().StringP("config", "c",
		path.Join(os.Getenv("HOME"), ".zebra.yaml"),
		"config file",
	)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	rootCmd.AddCommand(NewZebra())

	return rootCmd
}
