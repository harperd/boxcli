package boxclient

import (
	"encoding/json"
	"bytes"
	"github.com/hokaccha/go-prettyjson"
	"strings"
	"log"
)

func toInterface(s string) interface{} {
	var i map[string]interface{}

	if err := json.Unmarshal([]byte(s), &i); err != nil {
		log.Fatal(err)
	}

	return i
}

func toInterfaceArray(s string) []interface{} {
	var i []interface{}

	if err := json.Unmarshal([]byte(s), &i); err != nil {
		log.Fatal(err)
	}

	return i
}

/*
func ToString(i interface{}) (string, error) {
	var result string
	b, err := json.Marshal(i)

	if err == nil {
		result = string(b)
	}

	return result, err
}
*/

func formatJson(jsonString string, cfg *Config) string {
	var js string

	if cfg.Options.Unformatted || strings.Index(jsonString, "{") == -1 {
		js = jsonString
	} else if cfg.Options.Color {
		js = formatJsonColor(jsonString);
	} else {
		js = formatJsonMono(jsonString);
	}

	return js
}

func formatJsonMono(jsonString string) string {
	var byteBuf bytes.Buffer

	if err := json.Indent(&byteBuf, []byte(jsonString), "", "  "); err != nil {
		log.Fatal(err)
	}

	return byteBuf.String()
}

func formatJsonColor(js string) string {
	var buf []byte
	var err error
	c := string(js[0])

	if c == "[" {
		j := toInterfaceArray(js);

		if buf, err = prettyjson.Marshal(j); err != nil {
			log.Fatal(err)
		}
	} else if c == "{" {
		j := toInterface(js)

		if buf, err = prettyjson.Marshal(j); err != nil {
			log.Fatal(err)
		}
	}

	return string(buf)
}
