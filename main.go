package main

import (
	"fmt"
	"os"
	"errors"
	"strings"
	"./boxclient"
)

func showHelp() {
	fmt.Println("Usage: box [GET] [Resource] [Options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("\t-M\tmonochrome (don't colorize JSON)")
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
	opt.Address = os.Getenv("BOXENV")

	if opt.Address == "" {
		return nil, errors.New("BOXENV not set")
	}

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
	var result string
	var err error

	opt, err := getOptions(os.Args)

	if err == nil {
		s, err := boxclient.Execute(opt)

		if err == nil {
			if len(opt.Query) > 0 {
				s, err = boxclient.ApplyJsonQuery(s, opt)
			}

			if err == nil {
				s, err = boxclient.FormatJson(s, opt)
			}
		}

		result = s
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "{0}", err)
	} else {
		if len(result) > 0 {
			fmt.Println(result)
		}
	}
}
