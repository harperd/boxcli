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
const DB_FHIR string = "fhir"

func execute(cfg *Config) (string, string, error) {
	var err error
	var req *http.Request
	var jsonb []byte
	var jsons string = ""
	var message string

	req, err = createBoxRequest(cfg)

	if err == nil {
		jsonb, message, err = executeRequest(req)

		if err == nil {
			jsons = string(jsonb)

			if len(jsons) > 0 {
				err = checkErrors(jsons)
			}
		}
	}

	return jsons, message, err
}

func createBoxRequest(cfg *Config) (*http.Request, error) {
	var err error
	var req *http.Request
	var method = strings.ToUpper(cfg.Connection.Method)

	url, err := getBoxUrl(cfg)

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

func getBoxUrl(cfg *Config) (string, error) {
	var err error
	var url = os.Getenv(fmt.Sprintf("BOX_%s", strings.ToUpper(cfg.Connection.Box)))

	if len(url) == 0 {
		err = errors.New(fmt.Sprintf("Box %s not found.", cfg.Connection.Box))
	} else {
		url = fmt.Sprintf("%[1]s/%[2]s/%[3]s", url, cfg.Connection.Database, cfg.Options.Resource)

		if(cfg.Connection.Database == DB_FHIR) {
			if strings.Index(cfg.Options.Resource, "?") >= 0 {
				url += "&"
			} else {
				url += "?"
			}

			url += "_count=" + MAX_RESOURCES
		}
	}

	return url, err
}

func executeRequest(req *http.Request) ([]byte, string, error) {
	var err error
	var resp *http.Response
	var client = &http.Client{}
	var jsonb []byte
	var message string

	resp, err = client.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			message = resp.Status
		} else {
			jsonb, err = ioutil.ReadAll(resp.Body)
		}
	}

	return jsonb, message, err
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