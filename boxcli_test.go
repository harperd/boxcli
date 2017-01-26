package main

import (
	"testing"
	"github.com/harperd/boxcli/boxclient"
	"fmt"
	"errors"
)

func TestExecute(t *testing.T) {
	const ERRSTR string = "TestExecute(): %s"
	var result string
	args := []string{"", "testbox", "get", "fhir", "Patient"}
	cfg, err := boxclient.GetConfig(args)

	if err == nil {
		result, err = boxclient.Apply(cfg)
	}

	if err != nil {
		t.Errorf(ERRSTR, err)
	} else if len(result) == 0 {
		t.Errorf(ERRSTR, "No results")
	}
}

func TestIndexFirst(t *testing.T) {
	const ERRSTR = "TestIndex(): %s"
	var result string

	args := []string{"", "testbox", "get", "fhir", "Patient", "-i:0"}
	cfg, err := boxclient.GetConfig(args)

	if err == nil {
		result, err = boxclient.Apply(cfg)
	}

	if err == nil {
		// -- Get the count
		args = []string{"", "testbox", "get", "fhir", "Patient", "-c"}
		opt, err := boxclient.GetConfig(args)

		if err == nil {
			result, err = boxclient.Apply(opt)

			if err == nil {
				if result != "1" {
					err = errors.New(fmt.Sprintf("Expected only 1 resource. Received %s.", result))
				}
			}
		}
	}

	if err != nil {
		t.Errorf(ERRSTR, err)
	} else if len(result) == 0 {
		t.Errorf(ERRSTR, "No results")
	} else {
		fmt.Println(result)
	}
}

func TestIndexLast(t *testing.T) {
	const ERRSTR = "TestIndex(): %s"
	var result string

	args := []string{"", "testbox", "get", "fhir", "Patient", "-i:last"}
	cfg, err := boxclient.GetConfig(args)

	if err == nil {
		result, err = boxclient.Apply(cfg)
	}

	if err == nil {
		// -- Get the count
		args = []string{"", "testbox", "get", "fhir", "Patient", "-c"}
		opt, err := boxclient.GetConfig(args)

		if err == nil {
			result, err = boxclient.Apply(opt)

			if err == nil {
				if result != "1" {
					err = errors.New(fmt.Sprintf("Expected only 1 resource. Received %s.", result))
				}
			}
		}
	}

	if err != nil {
		t.Errorf(ERRSTR, err)
	} else if len(result) == 0 {
		t.Errorf(ERRSTR, "No results")
	} else {
		fmt.Println(result)
	}
}

// box testbox get fhir CareTeam .participant[].member.reference