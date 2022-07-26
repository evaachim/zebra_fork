package main

import (
	"fmt"
	"os"

	"github.com/project-safari/zebra/cli/show"
)

func main() {
	rootCmd := show.New()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
