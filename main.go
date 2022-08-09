package main

import (
	"fmt"
	"os"

	"git.ecd.axway.org/tjohnson/layer7/pkg/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
