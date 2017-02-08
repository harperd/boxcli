package boxclient

import (
	"github.com/jingweno/jqpipe-go"
	"encoding/json"
	"bytes"
	"strings"
	"strconv"
	"fmt"
	"os"
)

func evalJq(q string, js string) ([]json.RawMessage, error)  {
	var seq []json.RawMessage
	var err error
	//var result string

	seq, err = jq.Eval(js, q)

	if err != nil {
		fmt.Printf("ERROR: jq -> %s\n", q)
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}

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

func applyJsonQuery(s string, cfg *Config) (string, error) {
	var result string
	var err error
	var seq []json.RawMessage

	if doCompile(cfg) {
		seq, err = evalJq(compileQuery(s, cfg), s)

		if err == nil {
			result, err = formatOutput(seq, cfg)
		}
	} else if(cfg.Options.Count) {
		var i int = -1
		i, err = getResourceCount(s, cfg)

		if err == nil {
			result = strconv.Itoa(i)
		}
	} else {
		result, err = formatJson(s, cfg)
	}

	return result, err
}

func doCompile(cfg *Config) bool {
	return len(cfg.JQ.Custom) > 0 || len(cfg.Options.Index) > 0
}

func compileQuery(s string, cfg *Config) string {
	var q string
	var list = isBundle(s) || isArray(s)

	if len(cfg.Options.Index) > 0 && list {
		index, err := getIndex(s, cfg)

		if err == nil {
			q = strings.Replace(cfg.JQ.List.Index, "{index}", index, -1)

			if len(cfg.JQ.Custom) > 0 {
				q += "|" + cfg.JQ.Custom
			}
		}
	} else {
		if list {
			if cfg.JQ.List.Resources == cfg.JQ.Custom {
				q = cfg.JQ.Custom
			} else {
				q = fmt.Sprintf("%[1]s|%[2]s", cfg.JQ.List.Resources, cfg.JQ.Custom)
			}
		} else {
			q = cfg.JQ.Custom
		}
	}

	if len(q) > 0 {
		q = unquote(q)
		q = strings.Replace(q, "{", "\"", -1)
		q = strings.Replace(q, "}", "\"", -1)
	}

	return q
}

func getResourceCount(s string, cfg *Config) (int, error) {
	var count int = -1

	seq, err := evalJq(cfg.JQ.List.Count, s)

	if err == nil {
		count, err = strconv.Atoi(string(seq[0]))
	}

	return count, err
}

func getIndex(s string, cfg *Config) (string, error) {
	var index string = cfg.Options.Index
	var err error

	if strings.ToUpper(cfg.Options.Index) == "LAST" {
		var i, err = getResourceCount(s, cfg)

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

func formatOutput(seq []json.RawMessage, cfg *Config) (string, error) {
	var result string
	var err error

	if len(seq) > 1 {
		if stringValues(seq) {
			result, err = toSimpleList(seq, cfg)
		} else {
			result, err = toJsonArray(seq, cfg)
		}
	} else {
		result, err = formatJson(string(seq[0]), cfg)
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

func addValue(s string, cfg *Config) bool {
	return !cfg.Options.OmitNulls || (cfg.Options.OmitNulls && s != "null")
}

func toSimpleList(seq []json.RawMessage, cfg *Config) (string, error) {
	return writeList(seq, cfg, "", "", "\n")
}

func toJsonArray(seq []json.RawMessage, cfg *Config) (string, error) {
	return writeList(seq, cfg, "[", "]", ", ")
}

func writeList(seq []json.RawMessage, cfg *Config, open string, close string, sep string) (string, error) {
	var err error
	var s string
	var buf bytes.Buffer
	var indexes = make([]int, len(seq))

	// -- Do a priming read to get the indexes we need.
	for i := 0; i < len(seq); i++ {
		if addValue(string(seq[i]), cfg) {
			indexes = append(indexes, i)
		}
	}

	buf.WriteString(open)

	for i := 0; i<len(indexes); i++ {
		s = string(seq[indexes[i]])
		s, err = formatJson(s, cfg)

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