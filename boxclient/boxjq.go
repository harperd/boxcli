package boxclient

import (
	"github.com/jingweno/jqpipe-go"
	"encoding/json"
	"bytes"
	"strings"
	"strconv"
	"fmt"
)

func evalJq(q string, js string) ([]json.RawMessage, error)  {
	var seq []json.RawMessage
	var err error
	//var result string

	seq, err = jq.Eval(js, q)

	/*if err != nil {
		fmt.Printf("ERROR: jq -> %s\n", q)
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}*/

	return seq, err
}

/*
func showSummary(json string) string {
	summary := ""

	if(isBundle(json)) {
		jsonQuery(".entry[].resource.id", json)
	}

	return summary
}
*/

func applyJsonQuery(s string, opt *Options) (string, error) {
	var result string
	var err error
	var seq []json.RawMessage

	if doCompile(opt) {
		seq, err = evalJq(compileQuery(s, opt), s)

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
		result, err = formatJson(s, opt)
	}

	return result, err
}

func doCompile(opt *Options) bool {
	return len(opt.Query) > 0 || len(opt.Index) > 0
}

func compileQuery(s string, opt *Options) string {
	var q string
	var bundle = isBundle(s)
	var array = isArray(s)

	if len(opt.Index) > 0 && (bundle || array) {
		index, err := getIndex(s, opt)

		if err == nil {
			q = strings.Replace(opt.JsonIndex, "{index}", index, -1)

			if len(opt.Query) > 0 {
				q += "|" + opt.Query
			}
		}
	} else {
		if bundle || array {
			q = fmt.Sprintf("%[1]s|%[2]s", opt.JsonBase, opt.Query)
		} else {
			q = opt.Query
		}
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
	q := fmt.Sprintf("%s|length", opt.JsonBase)

	seq, err := evalJq(q, s)

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
				index = strconv.Itoa(i - 1)
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
		result, err = formatJson(string(seq[0]), opt)
	}

	return result, err
}

func isBundle(s string) bool {
	var bundle = false
	seq, _ := evalJq("select(.resourceType==\"Bundle\")|length", s)

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
	seq, _ := evalJq(".[]|length", s)

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
	return writeList(seq, opt, "", "", "\n")
}

func toJsonArray(seq []json.RawMessage, opt *Options) (string, error) {
	return writeList(seq, opt, "[", "]", ", ")
}

func writeList(seq []json.RawMessage, opt *Options, open string, close string, sep string) (string, error) {
	var err error
	var s string
	var buf bytes.Buffer
	var indexes = make([]int, len(seq))

	// -- Do a priming read to get the indexes we need.
	for i := 0; i < len(seq); i++ {
		if addValue(string(seq[i]), opt) {
			indexes = append(indexes, i)
		}
	}

	buf.WriteString(open)

	for i := 0; i<len(indexes); i++ {
		s = string(seq[indexes[i]])
		s, err = formatJson(s, opt)

		if err == nil {
			buf.WriteString(unquote(s))

			if i < len(indexes) - 1 {
				buf.WriteString(sep)
			}
		}
	}

	buf.WriteString(close)
	return buf.String(), err
}