package boxclient

import (
	"github.com/jingweno/jqpipe-go"
	"encoding/json"
	"bytes"
)

func ApplyJsonQuery(s string, opt *Options) (string, error) {
	var result string
	var err error

	if opt.Query != "" {
		q := opt.Query

		if isBundle(s) {
			q = ".entry[].resource|" + opt.Query
		}

		seq, err := jq.Eval(s, q)

		if err == nil {
			if len(seq) > 1 {
				result, err = toArray(seq, opt)
			} else {
				result, err = FormatJson(string(seq[0]), opt)
			}
		}
	} else {
		result, err = FormatJson(s, opt)
	}

	return result, err
}

func isBundle(s string) bool {
	var bundle = false
	seq, _ := jq.Eval(s, "select(.resourceType==\"Bundle\")")

	if len(seq) > 0 {
		bundle = true
	}

	return bundle
}

func toArray(seq []json.RawMessage, opt *Options) (string, error) {
	var err error
	var s string
	var buf bytes.Buffer

	buf.WriteString("[\n")

	for i := 0; i < len(seq); i++ {
		if err == nil {
			if !opt.OmitNulls || (opt.OmitNulls && string(seq[i]) != "null") {
				s, err = FormatJson(string(seq[i]), opt)

				if err == nil {
					buf.WriteString(s)

					if i < len(seq) - 1 {
						buf.WriteString(",")
					}

					buf.WriteString("\n")
				}
			}
		}
	}

	buf.WriteString("]\n")

	return buf.String(), err
}