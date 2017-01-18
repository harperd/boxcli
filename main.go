package main

import (
	"fmt"
	"os"
	"github.com/harperd/boxcli/boxclient"
)

func showHelp() {
	fmt.Println("Usage: box [box name] [get|put|post|delete] [doc|fhir] [resource] [options] <jq filter>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("\t-M\tmonochrome (don't colorize JSON)")
	fmt.Println("\t-u\tunformatted output")
	fmt.Println("\t-c\tget the count of the query results only")
	fmt.Println("\t-i:n\tget the resource at index n. Other value for n is 'last'.")
	os.Exit(0)
}

func main() {
	var err error
	var cfg *boxclient.Config
	var s string

	cfg, err = boxclient.GetConfig(os.Args)

	if err == nil {
		s, err = boxclient.Apply(cfg)
	} else {
		showHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	} else {
		if len(s) > 0 {
			fmt.Println(s)
		}
	}
}
