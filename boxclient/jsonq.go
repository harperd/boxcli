package boxclient

import "github.com/jingweno/jqpipe-go"

func ApplyJsonQuery(s string, opt *Options) (string, error) {
	var result string = s
	var err error

	if len(opt.Query) > 0 {
		seq, err := jq.Eval(s, opt.Query)

		if err == nil {
			result = string(seq[0])
		}
	}

	return result, err
}