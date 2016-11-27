package boxclient

import (
	"github.com/jingweno/jqpipe-go"
)

func ApplyJsonQuery(s string, opt *Options) (string, error) {
	var result string = s
	var err error

	seq, err := jq.Eval(s, opt.Query)

	if err == nil {
		for i := 0; i < len(seq); i++ {
			result += string(seq[i])
		}
	}

	return result, err
}