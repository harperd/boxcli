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

func createBoxRequest(cfg *Config) (*http.Request) {
	var method = strings.ToUpper(cfg.Connection.Method)

	url, auth := getBoxSettings(cfg)
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
	}

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

	return req
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

func getBoxSettings(cfg *Config) (string, string) {
	var url string
	var auth string
	var settings = os.Getenv(fmt.Sprintf("BOX_%s", strings.ToUpper(cfg.Connection.Box)))

	if len(settings) == 0 {
		log.Fatal(fmt.Sprintf("Box %s not found.", cfg.Connection.Box))
	} else {
		var tokens = strings.Split(settings, ";")

		if len(tokens) < 2 {
			log.Fatal("Invalid box settings")
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

			auth = base64.StdEncoding.EncodeToString([]byte(tokens[1]))
		}
	}

	return url, auth
}

func executeRequest(req *http.Request) ([]byte, string) {
	var client = &http.Client{}
	var jsonb []byte
	var message string
	var err error

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		message = resp.Status
	} else {
		jsonb, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}
	}

	return jsonb, message
}

func checkErrors(cfg *Config, js string) error {
	var err error = nil

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