package boxclient

func ApplyJsonQuery(jsonString string, opt *Options) (string, error) {
	var err error
	var result string
	/*
		if len(opt.Query) > 0 {
			parser, err := jsonql.NewStringQuery(jsonString)

			if err == nil {
				s, err := parser.Query(opt.Query)

				if err == nil {
					bytes, err := json.Marshal(s)

					if err == nil {
						result = string(bytes)
					}
				}
			}
		} else {*/
	result = jsonString
	//}

	return result, err
}