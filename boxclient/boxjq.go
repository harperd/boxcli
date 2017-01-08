package boxclient

import (
	"github.com/jingweno/jqpipe-go"
	"encoding/json"
	"bytes"
	"strings"
	"strconv"
)

func JQ(q string, js string) string  {
	var seq []json.RawMessage
	var err error
	var result string

	seq, err = jq.Eval(js, q)

	if err == nil {
		if len(seq) > 0 {
			result = string(seq[0])
		}
	}

	return result
}

func ShowSummary(json string) string {
	summary := ""

	if(isBundle(json)) {
		JQ(".entry[].resource.id", json)
	}

	return summary
}

func ApplyJsonQuery(s string, opt *Options) (string, error) {
	var result string
	var err error
	var seq []json.RawMessage

	if doCompile(opt) {
		seq, err = jq.Eval(s, compileQuery(s, opt))

		if err == nil {
			result, err = formatOutput(seq, opt)
		}
	} else if(opt.Count) {
		var i int = -1
		i, err = getResourceCount(s, opt)

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

	if len(opt.Index) > 0 && (isBundle(s) || isArray(s)) {
		index, err := getIndex(s, opt)

		if err == nil {
			q = strings.Replace(opt.JsonIndex, "{index}", index, -1)

			if len(opt.Query) > 0 {
				q += "|" + opt.Query
			}
		}
	} else {
		q = opt.Query
	}

	if len(q) > 0 {
		q = unquote(q)
		q = strings.Replace(q, "{", "\"", -1)
		q = strings.Replace(q, "}", "\"", -1)
	}

	return q
}

func getResourceCount(s string, opt *Options) (int, error) {
	var count int = -1

	seq, err := jq.Eval(s, opt.JsonBase + "|length")

	if err == nil {
		count, err = strconv.Atoi(string(seq[0]))
	}

	return count, err
}

func getIndex(s string, opt *Options) (string, error) {
	var index string = opt.Index
	var err error

	if strings.ToUpper(opt.Index) == "LAST" {
		var i, err = getResourceCount(s, opt)

		if err == nil {
			if err == nil && i > 0 {
				i--
				index = strconv.Itoa(i)
			}
		}
	}

	return index, err
}

func unquote(s string) string {
	s = strings.Trim(s, "")
	c := s[0]

	if c == '"' || c == '\'' {
		s = s[1:len(s) - 1]
	}

	return s
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
	seq, _ := jq.Eval(s, "select(.resourceType==\"Bundle\")|length")

	if len(seq) > 0 {
		i, _ := strconv.Atoi(string(seq[0]))

		if i > 0 {
			bundle = true
		}
	}

	return bundle
}

func isArray(s string) bool {
	var docArray = false
	seq, _ := jq.Eval(s, ".[]|length")

	if len(seq) > 0 {
		i, _ := strconv.Atoi(string(seq[0]))

		if i > 0 {
		docArray = true
		}
	}

	return docArray
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