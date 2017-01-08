package main

import (
	"testing"
	"github.com/harperd/boxcli/boxclient"
	"os"
	"strings"
	"fmt"
	"errors"
)

func TestExecute(t *testing.T) {
	const ERRSTR string = "TestExecute(): %s"
	opt := createOptions()

	opt.Database = "fhir"
	opt.Method = "get"
	opt.Resource = "Patient"

	s, err := boxclient.Execute(opt)

	if err == nil {
		err = checkErrors(s)
	}

	if err != nil {
		t.Errorf(ERRSTR, err)
	} else if len(s) == 0 {
		t.Errorf(ERRSTR, "No results")
	}
}

func TestFormat(t *testing.T) {
	const ERRSTR string = "TestFormat(): %s"
	const JSON string = "{ \"name\": \"john doe\", \"items\": { \"item1\": 1, \"item2\": 2 }, \"list\": [{\"listItem1\": 1}, { \"listItem2\": 2}] }"

	opt := createOptions()

	opt.Color = true
	opt.Unformatted = false

	f, err := boxclient.FormatJson(JSON, opt)

	if err != nil {
		t.Errorf(ERRSTR, err)
	} else if len(f) == 0 {
		t.Errorf(ERRSTR, "Format failed")
	} else if strings.Index(f, "\n") == -1 {
		t.Errorf(ERRSTR, "Format failed")
	} else if len(f) <= len(JSON) {
		t.Errorf(ERRSTR, "Format failed")
	}
}

func TestIndex(t *testing.T) {
	const ERRSTR = "TestIndex(): %s"

	opt := createOptions()

	opt.Database = "fhir"
	opt.Index = "0"
	opt.Method = "get"
	opt.Resource = "Patient"

	s, err := boxclient.Execute(opt)

	if err == nil {
		err = checkErrors(s)

		if err == nil {
			opt.Index = ""
			opt.Count = true

			s, err = boxclient.ApplyJsonQuery(s, opt)

			if err == nil {
				if s != "1" {
					err = errors.New("Expected only 1 resource. Received " + s + ".")
				}
			}
		}
	}

	if err != nil {
		t.Errorf(ERRSTR, err)
	} else if len(s) == 0 {
		t.Errorf(ERRSTR, "No results")
	} else {
		fmt.Println(s)
	}
}

func checkErrors(js string) error {
	var err error = nil

	var i interface{}
	i, err = boxclient.ToInterface(js)

	if err == nil {
		j := i.(map[string]interface{})
		if msg, ok := j["message"]; ok {
			err = errors.New(msg.(string))
		}
	}

	return err
}

func createOptions() *boxclient.Options {
	if (len(os.Getenv("BOXURL")) == 0) {
		os.Stderr.WriteString("BOXURL not set!");
	}

	return new(boxclient.Options)
}