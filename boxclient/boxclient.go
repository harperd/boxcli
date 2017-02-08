package boxclient

import (
	"net/http"
	"strings"
	"errors"
	"os"
	"io/ioutil"
	"fmt"
	"encoding/base64"
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
				err = checkErrors(cfg, jsons)
			}
		}
	}

	return jsons, message, err
}

func createBoxRequest(cfg *Config) (*http.Request, error) {
	var err error
	var req *http.Request
	var method = strings.ToUpper(cfg.Connection.Method)
	var url string
	var auth string

	url, auth, err = getBoxSettings(cfg)

	if err == nil {
		req, err = http.NewRequest(method, url, nil)

		if err == nil {
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic " + auth)

			/*
			var token string
			token, err = getJwt()

			if err == nil {
				fmt.Println(token)
				req.Header.Add("Authorization", "Bearer " + token)
			}
			*/

			if method == "POST" || method == "PUT" {
				req.Header.Add("Accepts", "application/json")
			}
		}
	}

	return req, err
}

/*
func getJwt() (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	claims := jwtToken.Claims.(jwt.MapClaims)

	claims["nickname"] = "boxcli"
	claims["iss"] = ""
	claims["sub"] = ""
	claims["aud"] = ""
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	mySigningKey := []byte("secret")

	return jwtToken.SignedString(mySigningKey)
}
*/

func getBoxSettings(cfg *Config) (string, string, error) {
	var err error
	var url string
	var auth string
	var settings = os.Getenv(fmt.Sprintf("BOX_%s", strings.ToUpper(cfg.Connection.Box)))

	if len(settings) == 0 {
		err = errors.New(fmt.Sprintf("Box %s not found.", cfg.Connection.Box))
	} else {
		var tokens = strings.Split(settings, ";")

		if len(tokens) < 2 {
			err = errors.New("Invalid box settings")
		} else {
			url = fmt.Sprintf("%[1]s/%[2]s/%[3]s", tokens[0], cfg.Connection.Database, cfg.Options.Resource)

			if (cfg.Connection.Database == DB_FHIR) {
				if strings.Index(cfg.Options.Resource, "?") >= 0 {
					url += "&"
				} else {
					url += "?"
				}

				url += "_count=" + MAX_RESOURCES
			}

			auth = getBasicAuthEncoded(tokens[1])
		}
	}

	return url, auth, err
}

func getBasicAuthEncoded (auth string) (string) {
	return base64.StdEncoding.EncodeToString([]byte(auth))
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

func checkErrors(cfg *Config, js string) error {
	var err error = nil

	// TODO: Fix for use with $documents
	if cfg.Connection.Box == "fhir" {
		var i interface{}
		i, err = toInterface(js)

		if err == nil {
			j := i.(map[string]interface{})
			if msg, ok := j["message"]; ok {
				err = errors.New(msg.(string))
			}
		}
	}

	return err
}