package boxclient

import (
	"encoding/json"
	"bytes"
	"github.com/hokaccha/go-prettyjson"
)

func formatJsonMono(jsonString string) (string, error) {
	var formatted string
	var byteBuf bytes.Buffer
	err := json.Indent(&byteBuf, []byte(jsonString), "", "  ")

	if err == nil {
		formatted = byteBuf.String()
	}

	return formatted, err
}

func formatJsonColor(jsonString string) (string, error) {
	var j map[string] interface{}
	json.Unmarshal([]byte(jsonString), &j)
	buf, err := prettyjson.Marshal(j)
	s := string(buf)
	return s, err
}

func FormatJson(jsonString string, opt *Options) (string, error) {
	var js string
	var err error

	if opt.Unformatted {
		js = jsonString
	} else {
		if (opt.Color) {
			js, err = formatJsonColor(jsonString)
		} else {
			js, err = formatJsonMono(jsonString)
		}
	}

	return js, err
}
