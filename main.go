package main

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
	"github.com/hokaccha/go-prettyjson"
	"github.com/elgs/jsonql"
	//"github.com/elgs/gojq"
	"bytes"
	"errors"
)

type Options struct {
	Address string
	Method string
	Resource string
	Color bool
	Query string
}

const MAX_RESOURCES string = "999999999"
const PROTOCOL string = "http"

func createBoxRequest(opt *Options) (*http.Request, error) {
	if len(opt.Address) == 0 {
		return nil, errors.New("Missing address")
	}

	req, err := http.NewRequest(strings.ToUpper(opt.Method),
		PROTOCOL + "://" + opt.Address + "/fhir/" + opt.Resource, nil)

	q := req.URL.Query()
	q.Add("_count", MAX_RESOURCES)
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accepts", "application/json")

	return req, nil
}

func executeRequest(req *http.Request) ([]byte, error) {
	var err error
	var resp *http.Response
	var client = &http.Client{}

	resp, err = client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func doRequest(opt *Options) (string, error) {
	var err error
	var req *http.Request
	var jsonb []byte

	req, err = createBoxRequest(opt)

	if err != nil {
		return "", err
	}

	jsonb, err = executeRequest(req)

	if err != nil {
		return "", err
	}

	return string(jsonb), err
}

func formatJsonMono(jsonString string /*jsonb []byte*/) (string, error) {
	var byteBuf bytes.Buffer
	err := json.Indent(&byteBuf, []byte(jsonString), "", "  ")

	if err != nil {
		return "", err
	}

	return byteBuf.String(), nil
}

func formatJsonColor(jsonString string) (string, error) {
	var j map[string] interface{}
	json.Unmarshal([]byte(jsonString), &j)
	buf, err := prettyjson.Marshal(j)
	s := string(buf)
	return s, err
}

func formatJson(jsonString string, opt *Options) (string, error) {
	var js string
	var err error

	if(opt.Color) {
		js, err = formatJsonColor(jsonString)
	} else {
		js, err = formatJsonMono(jsonString)
	}

	return js, err
}

func applyJsonQuery(jsonString string, opt *Options) (string, error) {
	var err error
	var result string

	if len(opt.Query) > 0 {
		parser, err := jsonql.NewStringQuery(jsonString)

		if err == nil {
			s, err := parser.Query(opt.Query)

			if err == nil {
				bytes, err := json.Marshal(s)

				if err == nil {
					result = string(bytes)
				}
			}
		}
	} else {
		result = jsonString
	}

	return result, err
}

func showHelp() {
	fmt.Println("Usage: box [GET] [Resource] [Options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("\t-M\tmonochrome (don't colorize JSON)")
	os.Exit(0)
}

func processArg(arg string, opt *Options) {
	for c := 1; c < len(arg); c++ {
		arg := string(arg[c])

		if arg == "M" {
			opt.Color = false;
		}
	}
}

func processArgs(args []string, opt *Options) {
	for c := 0; c < len(args); c++ {
		arg := args[c]
		if strings.Index(arg, "-") == 0 {
			processArg(arg, opt)
		} else {
			opt.Query = arg
		}
	}
}

func getOptions(args []string) (*Options, error) {
	opt := new(Options);
	opt.Color = true;
	opt.Address = os.Getenv("BOX_ENV")

	if opt.Address == "" {
		return nil, errors.New("BOX_ENV not set")
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
		s, err := doRequest(opt)

		if err == nil {
			if len(opt.Query) > 0 {
				s, err = applyJsonQuery(s, opt)
			}

			if err == nil {
				s, err = formatJson(s, opt)
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
