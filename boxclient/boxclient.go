package boxclient

import (
	"net/http"
	"strings"
	"errors"
	"os"
	"io/ioutil"
	"fmt"
)

const MAX_RESOURCES string = "999999999"

func execute(opt *Options) (string, error) {
	var err error
	var req *http.Request
	var jsonb []byte
	var jsons string = ""

	req, err = createBoxRequest(opt)

	if err == nil {
		jsonb, err = executeRequest(req)

		if err == nil {
			jsons = string(jsonb)
			err = checkErrors(jsons)
		}
	}

	return jsons, err
}

func createBoxRequest(opt *Options) (*http.Request, error) {
	var err error
	var req *http.Request
	var method = strings.ToUpper(opt.Method)

	url, err := getBoxUrl(opt)

	if err == nil {
		req, err = http.NewRequest(method, url, nil)

		if err == nil {
			req.Header.Add("Content-Type", "application/json")

			if method == "POST" || method == "PUT" {
				req.Header.Add("Accepts", "application/json")
			}
		}
	}

	return req, err
}

func getBoxUrl(opt *Options) (string, error) {
	var err error
	var url = os.Getenv(fmt.Sprintf("BOX_%s", strings.ToUpper(opt.Box)))

	if len(url) == 0 {
		err = errors.New(fmt.Sprintf("Box %s not found.", opt.Box))
	} else {
		url = fmt.Sprintf("%[1]s/%[2]s/%[3]s", url, opt.Database, opt.Resource)

		if(opt.Database == "fhir") {
			if strings.Index(opt.Resource, "?") >= 0 {
				url += "&"
			} else {
				url += "?"
			}

			url += "_count=" + MAX_RESOURCES
		}
	}

	return url, err
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

func checkErrors(js string) error {
	var err error = nil

	var i interface{}
	i, err = toInterface(js)

	if err == nil {
		j := i.(map[string]interface{})
		if msg, ok := j["message"]; ok {
			err = errors.New(msg.(string))
		}
	}

	return err
}