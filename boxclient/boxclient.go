package boxclient

import (
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
)

type Options struct {
	Address string
	Method string
	Resource string
	Color bool
	Unformatted bool
	Query string
}

const MAX_RESOURCES string = "999999999"
const PROTOCOL string = "http"

func createBoxRequest(opt *Options) (*http.Request, error) {
	if len(opt.Address) == 0 {
		return nil, errors.New("Missing address")
	}

	req, err := http.NewRequest(strings.ToUpper(opt.Method),
		PROTOCOL + "://" + opt.Address + "/fhir/" + opt.Resource + "?_count=" + MAX_RESOURCES, nil)

	if err == nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accepts", "application/json")
	}

	return req, nil
}

func executeRequest(req *http.Request) ([]byte, error) {
	var err error
	var resp *http.Response
	var client = &http.Client{}
	var jsonb []byte

	resp, err = client.Do(req)

	if err == nil {
		defer resp.Body.Close()
		jsonb, err = ioutil.ReadAll(resp.Body)
	}

	return jsonb, err
}

func Execute(opt *Options) (string, error) {
	var err error
	var req *http.Request
	var jsonb []byte

	req, err = createBoxRequest(opt)

	if err == nil {
		jsonb, err = executeRequest(req)
	}

	return string(jsonb), err
}
