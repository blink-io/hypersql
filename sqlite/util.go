package sqlite

func isSynonym(params map[string]string, keys ...string) bool {
	if len(keys) > 1 {
		var c = 0
		for _, key := range keys {
			if _, ok := params[key]; ok {
				c++
				if c > 1 {
					return true
				}
			}
		}
	}
	return false
}

func ifSynonym(params map[string]string, then func(...string) error, keys ...string) error {
	if isSynonym(params, keys...) {
		return then(keys...)
	}
	return nil
}
