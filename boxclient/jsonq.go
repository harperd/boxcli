package boxclient

func ApplyJsonQuery(s string, opt *Options) (string, error) {
	var result string = s
	var err error

	/*
	if len(opt.Query) > 0 {
		result, err = jq.Apply(opt.Query, ToInterface(s))

		if err == nil {

		}
	}
	*/

	return result, err
}