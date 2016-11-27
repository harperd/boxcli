package boxclient

import (
	"github.com/jingweno/jqpipe-go"
	"encoding/json"
)

func ApplyJsonQuery(s string, opt *Options) (string, error) {
	var result string
	var err error

	if opt.Query != "" {
		seq, err := jq.Eval(s, opt.Query)

		if err == nil {
			if len(seq) > 1 {
				result, err = toArray(seq, opt)
			} else {
				result, err = FormatJson(string(seq[0]), opt)
			}
		}
	} else {
		result = s
	}

	return result, err
}

func toArray(seq []json.RawMessage, opt *Options) (string, error) {
	var err error
	var s string
	result := "[\n"

	for i := 0; i < len(seq); i++ {
		if err == nil {
			s, err = FormatJson(string(seq[i]), opt)

			if err == nil {
				result += s

				if i < len(seq) - 1 {
					result += ","
				}

				result += "\n"
			}
		}
	}

	result += "]"

	return result, err
}