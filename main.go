package main

import (
	"fmt"
	"os"
	"github.com/harperd/boxcli/boxclient"
)

/*
box test get fhir Patient '.entry[].resource|select(.id=="db327e81-cef5-4f1b-b20a-0f2332b02584")'
ERROR: jq: error (at <stdin>:0): Cannot iterate over null (null)
 */

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
	var json string
	var message string

	cfg, err = boxclient.GetConfig(os.Args)

	if err == nil {
		json, message, err = boxclient.Apply(cfg)
	} else {
		showHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	} else if len(message) > 0 {
		fmt.Println(message)
	} else if len(json) > 0 {
		fmt.Println(json)
	}
}
