package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/harperd/boxcli/boxclient"
)

func showHelp() {
	fmt.Println("Usage: box [get|put|post|delete] [resource] [options] <jq filter>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("\t-M\tmonochrome (don't colorize JSON)")
	fmt.Println("\t-u\tunformatted output")
	os.Exit(0)
}

func processArg(arg string, opt *boxclient.Options) {
	for c := 1; c < len(arg); c++ {
		arg := string(arg[c])

		if arg == "M" {
			opt.Color = false
		} else if arg == "u" {
			opt.Unformatted = true
		}
	}
}

func processArgs(args []string, opt *boxclient.Options) {
	for c := 0; c < len(args); c++ {
		arg := args[c]
		if strings.Index(arg, "-") == 0 {
			processArg(arg, opt)
		} else {
			opt.Query = arg
		}
	}
}

func getOptions(args []string) (*boxclient.Options, error) {
	opt := new(boxclient.Options);
	opt.Color = true;
	opt.Unformatted = false;

	if len(args) >= 3 {
		opt.Method = args[1]
		opt.Resource = args[2]

		if len(args) > 3 {
			processArgs(args[3:], opt)
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
