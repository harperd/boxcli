package boxclient

import (
	"github.com/jingweno/jqpipe-go"
	"encoding/json"
	"bytes"
	"strings"
	"strconv"
	"fmt"
	"log"
)

func evalJq(q string, js string) []json.RawMessage {
	if seq, err := jq.Eval(js, q); err != nil {
		fmt.Printf("ERROR: jq -> %s\n", q)
		log.Fatal(err)
		return nil
	} else {
		return seq
	}
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

func applyJsonQuery(s string, cfg *Config) string {
	var result string
	var seq []json.RawMessage

	if doCompile(cfg) {
		seq = evalJq(compileQuery(s, cfg), s)
		result = formatOutput(seq, cfg)
	} else if (cfg.Options.Count) {
		i := getResourceCount(s, cfg)
		result = strconv.Itoa(i)
	} else {
		result = formatJson(s, cfg)
	}

	return result
}

func doCompile(cfg *Config) bool {
	return len(cfg.JQ.Custom) > 0 || len(cfg.Options.Index) > 0
}

func compileQuery(s string, cfg *Config) string {
	var q string
	var list = isBundle(s) || isArray(s)

	if len(cfg.Options.Index) > 0 && list {
		index := getIndex(s, cfg)
		q = strings.Replace(cfg.JQ.List.Index, "{index}", index, -1)

		if len(cfg.JQ.Custom) > 0 {
			q += "|" + cfg.JQ.Custom
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

func getResourceCount(s string, cfg *Config) int {
	var i int
	var err error

	seq := evalJq(cfg.JQ.List.Count, s)
	if i, err = strconv.Atoi(string(seq[0])); err != nil {
		log.Fatal(err)
	}

	return i
}

func getIndex(s string, cfg *Config) string {
	var index string = cfg.Options.Index

	if strings.ToUpper(cfg.Options.Index) == "LAST" {
		if i := getResourceCount(s, cfg); i > 0 {
			index = strconv.Itoa(i - 1)
		}
	}

	return index
}

func unquote(s string) string {
	s = strings.Trim(s, "")
	c := s[0]

	if c == '"' || c == '\'' {
		s = s[1:len(s) - 1]
	}

	return s
}

func formatOutput(seq []json.RawMessage, cfg *Config) string {
	var result string

	if len(seq) > 1 {
		if stringValues(seq) {
			result = toSimpleList(seq, cfg)
		} else {
			result = toJsonArray(seq, cfg)
		}
	} else {
		result = formatJson(string(seq[0]), cfg)
	}

	return result
}

func isBundle(s string) bool {
	var bundle = false
	if seq := evalJq("select(.resourceType==\"Bundle\")|length", s); len(seq) > 0 {
		if i, err := strconv.Atoi(string(seq[0])); err != nil {
			log.Fatal(err)
		} else {
			if i > 0 {
				bundle = true
			}
		}
	}

	return bundle
}

func isArray(s string) bool {
	var docArray = false
	if seq := evalJq(".[]|length", s); len(seq) > 0 {
		if i, err := strconv.Atoi(string(seq[0])); err != nil {
			log.Fatal(err)
		} else {
			if i > 0 {
				docArray = true
			}
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

func toSimpleList(seq []json.RawMessage, cfg *Config) string {
	return writeList(seq, cfg, "", "", "\n")
}

func toJsonArray(seq []json.RawMessage, cfg *Config) string {
	return writeList(seq, cfg, "[", "]", ", ")
}

func writeList(seq []json.RawMessage, cfg *Config, open string, close string, sep string) string {
	var buf bytes.Buffer
	var indexes = make([]int, len(seq))

	// -- Do a priming read to get the indexes we need.
	for i := 0; i < len(seq); i++ {
		if addValue(string(seq[i]), cfg) {
			indexes = append(indexes, i)
		}
	}

	buf.WriteString(open)

	for i := 0; i < len(indexes); i++ {
		s := formatJson(string(seq[indexes[i]]), cfg)

		buf.WriteString(unquote(s))

		if i < len(indexes) - 1 {
			buf.WriteString(sep)
		}
	}

	buf.WriteString(close)
	return buf.String()
}