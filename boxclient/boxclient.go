package boxclient

import (
	"net/http"
	"strings"
	"errors"
	"os"
	"io/ioutil"
	"fmt"
	"encoding/base64"
	"log"
)

const (
	MAX_RESOURCES string = "999999999"
	DB_FHIR string = "fhir"
)

func execute(cfg *Config) (string, string, error) {
	var err error
	var jsons string

	req := createBoxRequest(cfg)
	jsonb, message := executeRequest(req)

	if jsons = string(jsonb); len(jsons) > 0 {
		err = checkErrors(cfg, jsons)
	}

	return jsons, message, err
}

func createBoxRequest(cfg *Config) *http.Request {
	settings := getBoxSettings(cfg)
	return createRequest(cfg, settings)
}

func createRequest(cfg *Config, settings string) *http.Request {
	req, err := http.NewRequest(
		strings.ToUpper(cfg.Connection.Method),
		getUrl(settings, cfg), nil)

	if err != nil {
		log.Fatal(err)
	}

	setHeader(req, cfg, settings)

	return req
}

func setHeader(req *http.Request, cfg *Config, settings string) {
	var method = strings.ToUpper(cfg.Connection.Method)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic " + getAuth(settings))

	if method == "POST" || method == "PUT" {
		req.Header.Add("Accepts", "application/json")
	}
}

func getUrl(settings string, cfg *Config) string {
	var tokens = strings.Split(settings, ";")
	url := fmt.Sprintf("%[1]s/%[2]s/%[3]s", tokens[0], cfg.Connection.Database, cfg.Options.Resource)

	if (cfg.Connection.Database == DB_FHIR) {
		if strings.Index(cfg.Options.Resource, "?") >= 0 {
			url += "&"
		} else {
			url += "?"
		}

		url += "_count=" + MAX_RESOURCES
	}

	return url
}

func getAuth(settings string) string {
	var tokens = strings.Split(settings, ";")
	return base64.StdEncoding.EncodeToString([]byte(tokens[1]))
}

func getBoxSettings(cfg *Config) string {
	var settings = os.Getenv(fmt.Sprintf("BOX_%s", strings.ToUpper(cfg.Connection.Box)))

	if len(settings) == 0 {
		log.Fatal(fmt.Sprintf("Box %s not found.", cfg.Connection.Box))
	} else {
		var tokens = strings.Split(settings, ";")

		if len(tokens) < 2 {
			log.Fatal("Invalid box settings")
		}
	}

	return settings
}

func executeRequest(req *http.Request) ([]byte, string) {
	var client = &http.Client{}
	var jsonb []byte
	var message string

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			message = resp.Status
		} else {
			var err error

			jsonb, err = ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal(err)
	}

	return jsonb, message
}

func checkErrors(cfg *Config, js string) error {
	var err error

	// TODO: Fix for use with $documents
	if cfg.Connection.Box == "fhir" {
		i := toInterface(js)

		j := i.(map[string]interface{})

		if msg, ok := j["message"]; ok {
			err = errors.New(msg.(string))
		}
	}

	return err
}