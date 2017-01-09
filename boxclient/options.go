package boxclient

import (
	"strings"
	"errors"
)

// Runtime options for Box Client
// Method: GET, PUT, POST, DELETE, etc.
// Database: FHIR or DOC
// JsonBase: Base of JSON list for resources or documents (i.e. .entry or .[])
// Resource: A valid FHIR resource
// Color: If true, JSON output is syntax highlighted
// Unformatted: If true, JSON output is not formatted
// Count: If true, only the count of the results are returned
// Index: If true, only the resource at the specified index is returned
// Query: The JSON query to apply
type Options struct {
	Address string
	Method string
	Database string
	JsonBase string
	JsonIndex string
	Resource string
	Color bool
	Unformatted bool
	OmitNulls bool
	Count bool
	Index string
	Query string
}

func ApplyOptions(opt *Options) (string, error) {
	var err error
	var s string

	if err == nil {
		s, err = execute(opt)

		if err == nil {
			s, err = applyJsonQuery(s, opt)
		}
	}

	return s, err
}

// Set the Options structure by processing the provided command line
// arguments. This structure is used by the boxclient package to process
// Aidbox requests.
func GetOptions(args []string) (*Options, error) {
	var err error
	opt := new(Options)

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
			opt.JsonIndex = ".entry[{index}].resource"
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
		err = errors.New("Invalid options")
	}

	return opt, err
}

// Processes all command line arguments. A command line argument is proceeded
// by a dash ('-'). Any argument that is not proceeded by a dash is assumed to
// be a JQ query. The first two arguments, method and resource, are required and not
// processed by this function.
func processArgs(args []string, opt *Options) {
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

// Processes all command line arguments. A command line argument is proceeded
// by a dash ('-'). Acceptable arguments can be -Mcu or -M -c -u individually.
func processArg(arg string, opt *Options) {
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