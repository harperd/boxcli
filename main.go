package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/harperd/boxcli/boxclient"
)

func showHelp() {
	fmt.Println("Usage: box [get|put|post|delete] [doc|fhir] [resource] [options] <jq filter>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("\t-M\tmonochrome (don't colorize JSON)")
	fmt.Println("\t-u\tunformatted output")
	fmt.Println("\t-c\tget the count of the query results only")
	fmt.Println("\t-i:n\tget the resource at index n. Other value for n is 'last'.")
	os.Exit(0)
}

// Processes all command line arguments. A command line argument is proceeded
// by a dash ('-'). Acceptable arguments can be -Mcu or -M -c -u individually.
func processArg(arg string, opt *boxclient.Options) {
	if strings.Index(arg, "-i:") > -1 {
		opt.Index = strings.Split(arg, ":")[1]
	} else {
		for c := 1; c < len(arg); c++ {
			arg := string(arg[c])

			if arg == "M" {
				opt.Color = false
			} else if arg == "u" {
				opt.Unformatted = true
			} else if arg == "c" {
				opt.Count = true
			}
		}
	}
}

// Processes all command line arguments. A command line argument is proceeded
// by a dash ('-'). Any argument that is not proceeded by a dash is assumed to
// be a JQ query. The first two arguments, method and resource, are required and not
// processed by this function.
func processArgs(args []string, opt *boxclient.Options) {
	for c := 0; c < len(args); c++ {
		arg := args[c]
		if strings.Index(arg, "-") == 0 {
			processArg(arg, opt)
		} else {
			if len(arg) > 0 {
				opt.Query = arg
			}
		}
	}
}

// Set the Options structure by processing the provided command line
// arguments. This structure is used by the boxclient package to process
// Aidbox requests.
func getOptions(args []string) (*boxclient.Options, error) {
	opt := new(boxclient.Options)

	// -- Option defaults
	opt.Database = "fhir"
	opt.Color = true
	opt.Unformatted = false
	opt.OmitNulls = true
	opt.Count = false
	opt.Index = ""

	if len(args) >= 4 {
		opt.Method = args[1]

		if(strings.ToLower(args[2]) == "fhir") {
			opt.Database = "fhir";
			opt.JsonBase = ".entry"
			opt.JsonIndex = "entry[{index}].resource"
		} else if (strings.ToLower(args[2]) == "doc") {
			opt.Database = "$documents"
			opt.JsonBase = ".[]"
			opt.JsonIndex = ".[{index}]"
		}

		opt.Resource = args[3]

		if len(args) > 4 {
			processArgs(args[4:], opt)
		}
	} else {
		showHelp();
	}

	return opt, nil
}

func main() {
	var err error
	var opt *boxclient.Options
	var s string

	opt, err = getOptions(os.Args)

	if err == nil {
		s, err = boxclient.Execute(opt)

		if err == nil {
			s, err = boxclient.ApplyJsonQuery(s, opt)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	} else {
		if len(s) > 0 {
			fmt.Println(s)
		}
	}
}
