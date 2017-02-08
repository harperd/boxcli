package boxclient

import (
	"strings"
	"errors"
	"fmt"
)

const BOX_IDX = 1
const METHOD_IDX = 2
const DB_IDX = 3
const RESOURCE_IDX = 4
const JQ_IDX =  5

type Config struct {
	BasicAuth struct {
	                Username string
					Password string
	              }
	Connection struct {
		           Box      string
		           Method   string
		           Database string
	           }
	Options    struct {
		           Resource    string
		           Index       string
		           Color       bool
		           Unformatted bool
		           OmitNulls   bool
		           Count       bool
	           }
	JQ         struct {
		           Custom   string
		           List     struct {
			                    Count     string
			                    Index     string
			                    Resources string
		                    }
		           Resource struct {
			                    Name string
		                    }
	           }
}

func Apply(cfg *Config) (string, string, error) {
	var err error
	var json string
	var message string

	if err == nil {
		json, message, err = execute(cfg)

		if err == nil {
			if strings.ToUpper(cfg.Connection.Method) == "DELETE" {
				message = fmt.Sprintf("%s deleted.", cfg.Options.Resource)
			} else {
				if len(json) > 0 {
					json, err = applyJsonQuery(json, cfg)
				}
			}
		}
	}

	return json, message, err
}

func GetConfig(args []string) (*Config, error) {
	var err error
	cfg := new(Config)

	// -- Configuration defaults
	cfg.Connection.Database = ""
	cfg.Options.Color = true
	cfg.Options.Unformatted = false
	cfg.Options.OmitNulls = true
	cfg.Options.Count = false
	cfg.Options.Index = ""

	if len(args) >= 5 {
		cfg.Connection.Box = args[BOX_IDX]
		cfg.Connection.Method = args[METHOD_IDX]

		if(strings.ToLower(args[DB_IDX]) == "fhir") {
			cfg.Connection.Database = "fhir";
			cfg.JQ.List.Resources = ".entry[].resource"
			cfg.JQ.List.Index = ".entry[{index}].resource"
			cfg.JQ.List.Count = ".entry|length"
		} else if (strings.ToLower(args[DB_IDX]) == "doc") {
			cfg.Connection.Database = "$documents"
			cfg.JQ.List.Resources = ".[]"
			cfg.JQ.List.Index = ".[{index}]"
			cfg.JQ.List.Count = ".|length"
		}

		cfg.Options.Resource = args[RESOURCE_IDX]

		if len(args) > 5 {
			processArgs(args[JQ_IDX:], cfg)
		}
	} else {
		err = errors.New("Invalid options")
	}

	return cfg, err
}

func processArgs(args []string, opt *Config) {
	for c := 0; c < len(args); c++ {
		arg := args[c]
		if strings.Index(arg, "-") == 0 {
			processArg(arg, opt)
		} else {
			if len(arg) > 0 {
				opt.JQ.Custom = arg
			}
		}
	}
}

func processArg(arg string, opt *Config) {
	if strings.Index(arg, "-i:") > -1 {
		opt.Options.Index = strings.Split(arg, ":")[1]
	} else {
		for c := 1; c < len(arg); c++ {
			arg := string(arg[c])

			if arg == "M" {
				opt.Options.Color = false
			} else if arg == "u" {
				opt.Options.Unformatted = true
			} else if arg == "c" {
				opt.Options.Count = true
			}
		}
	}
}