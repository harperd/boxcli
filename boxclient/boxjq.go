package boxclient

import (
	"github.com/jingweno/jqpipe-go"
	"encoding/json"
	"bytes"
	"strings"
	"regexp"
)

func ApplyJsonQuery(s string, opt *Options) (string, error) {
	var result string
	var err error

	if opt.Query != "" || opt.Count || opt.Index {
		seq, err := jq.Eval(s, compileQuery(s, opt))

		if err == nil {
			result, err = formatOutput(seq, opt)
		}
	} else {
		result, err = FormatJson(s, opt)
	}

	return result, err
}

func compileQuery(s string, opt *Options) string {
	q := opt.Query

	if isBundle(s) {
		q = ".entry[" + opt.Index + "].resource|" + opt.Query
	}

	if opt.Count && len(q) > 0 {
		s += "|length"
	}

	return q
}

func unquote(s string) string {
	reg := regexp.MustCompile(`"([^"]*)"`)
	return reg.ReplaceAllString(s, "${1}")
}

func formatOutput(seq []json.RawMessage, opt *Options) (string, error) {
	var result string
	var err error

	if len(seq) > 1 {
		if stringValues(seq) {
			result, err = toSimpleList(seq, opt)
		} else {
			result, err = toJsonArray(seq, opt)
		}
	} else {
		result, err = FormatJson(string(seq[0]), opt)
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

func stringValues(seq []json.RawMessage) bool {
	return strings.Index(string(seq[0]), "\"") == 0
}

func addValue(s string, opt *Options) bool {
	return !opt.OmitNulls || (opt.OmitNulls && s != "null")
}

func toSimpleList(seq []json.RawMessage, opt *Options) (string, error) {
	var err error
	var buf bytes.Buffer

	for i := 0; i < len(seq); i++ {
		s := string(seq[i])

		if addValue(s, opt) {
			if i < len(seq) - 1 {
				buf.WriteString(unquote(s))

				if i < len(seq) - 2 {
					buf.WriteString("\n")
				}
			}
		}
	}

	return buf.String(), err
}

func toJsonArray(seq []json.RawMessage, opt *Options) (string, error) {
	var err error
	var s string
	var buf bytes.Buffer

	buf.WriteString("[")

	for i := 0; i < len(seq); i++ {
		s = string(seq[i])

		if addValue(s, opt) {
			s, err = FormatJson(s, opt)

			if err == nil {
				buf.WriteString(s)

				if i < len(seq) - 1 {
					buf.WriteString(",\n")
				}
			}
		}
	}

	buf.WriteString("]")
	return buf.String(), err
}