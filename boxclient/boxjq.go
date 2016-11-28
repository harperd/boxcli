package boxclient

import (
	"github.com/jingweno/jqpipe-go"
	"encoding/json"
	"bytes"
	"strings"
	"regexp"
	"strconv"
)

func ApplyJsonQuery(s string, opt *Options) (string, error) {
	var result string
	var err error

	if doCompile(opt) {
		seq, err := jq.Eval(s, compileQuery(s, opt))

		if err == nil {
			result, err = formatOutput(seq, opt)
		}
	} else if(opt.Count) {
		var i int = -1
		i, err = getResourceCount(s)

		if err == nil {
			result = strconv.Itoa(i)
		}
	} else {
		result, err = FormatJson(s, opt)
	}

	return result, err
}

func doCompile(opt *Options) bool {
	return len(opt.Query) > 0 || len(opt.Index) > 0
}

func compileQuery(s string, opt *Options) string {
	var q string

	if isBundle(s) {
		index, err := getIndex(s, opt)

		if err == nil {
			q = ".entry[" + index + "].resource"

			if len(opt.Query) > 0 {
				q += "|" + opt.Query
			}
		}
	} else {
		q = opt.Query
	}

	return q
}

func getResourceCount(s string) (int, error) {
	var count int = -1
	seq, err := jq.Eval(s, ".entry|length")

	if err == nil {
		count, err = strconv.Atoi(string(seq[0]))
	}

	return count, err
}

func getIndex(s string, opt *Options) (string, error) {
	var index string = opt.Index
	var err error
	var seq []json.RawMessage

	if strings.ToUpper(opt.Index) == "LAST" {
		seq, err = jq.Eval(s, ".entry|length")

		if err == nil {
			var i int = -1
			i, err = strconv.Atoi(string(seq[0]))

			if err == nil && i > 0 {
				i--
				index = strconv.Itoa(i)
			}
		}
	}

	return index, err
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