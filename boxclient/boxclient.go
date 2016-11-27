package boxclient

import (
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
	"os"
)

const MAX_RESOURCES string = "999999999"

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

func createBoxRequest(opt *Options) (*http.Request, error) {
	var url = os.Getenv("BOXURL")

	if len(url) == 0 {
		return nil, errors.New("BOXURL not set")
	}

	var method = strings.ToUpper(opt.Method)

	req, err := http.NewRequest(method, url + "/fhir/" + opt.Resource + "?_count=" + MAX_RESOURCES, nil)

	if err == nil {
		req.Header.Add("Content-Type", "application/json")

		if method == "POST" || method == "PUT" {
			req.Header.Add("Accepts", "application/json")
		}
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
