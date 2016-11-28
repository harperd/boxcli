package main

import (
	"testing"
	"github.com/harperd/boxcli/boxclient"
	"os"
	"strings"
)

func TestExecute(t *testing.T) {
	const ERRSTR string = "TestExecute(): %s"
	opt := testSetup()
	s, err := boxclient.Execute(opt)

	if err != nil {
		t.Errorf(ERRSTR, err)
	} else if len(s) == 0 {
		t.Errorf(ERRSTR, "No results")
	}
}

func TestFormat(t *testing.T) {
	const ERRSTR string = "TestFormat(): %s"
	const JSON string = "{ \"name\": \"john doe\", \"items\": { \"item1\": 1, \"item2\": 2 }, \"list\": [{\"listItem1\": 1}, { \"listItem2\": 2}] }"

	opt := testSetup()

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

func testSetup() *boxclient.Options {
	os.Setenv("BOXURL", "http://narus.aidbox.master.narus.aidbox.io")

	opt := new(boxclient.Options)
	opt.Color = true
	opt.Unformatted = false
	opt.OmitNulls = true
	opt.Count = false
	opt.Index = ""
	opt.Method = "get"
	opt.Resource = "Patient"

	return opt
}