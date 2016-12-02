package boxclient

import (
	"encoding/json"
	"bytes"
	"github.com/hokaccha/go-prettyjson"
	"strings"
	"fmt"
)

func ToInterface(s string) (interface{}, error) {
	var i map[string]interface{}
	err := json.Unmarshal([]byte(s), &i)
	return i, err
}

func ToInterfaceArray(s string) ([]interface{}, error) {
	var i []interface{}
	err := json.Unmarshal([]byte(s), &i)

	return i, err
}

func ToString(i interface{}) (string, error) {
	var result string
	b, err := json.Marshal(i)

	if err == nil {
		result = string(b)
	}

	return result, err
}

func FormatJson(jsonString string, opt *Options) (string, error) {
	var js string
	var err error

	if opt.Unformatted || strings.Index(jsonString, "{") == -1 {
		js = jsonString
	} else if opt.Color {
		js, err = formatJsonColor(jsonString, opt)
	} else {
		js, err = formatJsonMono(jsonString)
	}

	return js, err
}

func formatJsonMono(jsonString string) (string, error) {
	var formatted string
	var byteBuf bytes.Buffer
	err := json.Indent(&byteBuf, []byte(jsonString), "", "  ")

	if err == nil {
		formatted = byteBuf.String()
	}

	return formatted, err
}

func formatJsonColor(js string, opt *Options) (string, error) {
	var s string
	var err error

	c := string(js[0])

	if c == "[" {
		var j []interface{}
		j, err = ToInterfaceArray(js);

		if err == nil {
			buf, err := prettyjson.Marshal(j)

			if err == nil {
				s = string(buf)
				fmt.Println(len(buf))
			}
		}
	} else if c == "{" {
		var j interface{}
		j, err = ToInterface(js)

		if err == nil {
			buf, err := prettyjson.Marshal(j)

			if err == nil {
				s = string(buf)
			}
		}
	}

	return s, err
}
