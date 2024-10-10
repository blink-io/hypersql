package hypersql

type ConfigParams map[string]string

func (p ConfigParams) Get(keys ...string) string {
	if len(keys) > 0 {
		for _, key := range keys {
			if v := p[key]; len(v) > 0 {
				return v
			}
		}
	}
	return ""
}

func (p ConfigParams) Exists(key string) bool {
	_, ok := p[key]
	return ok
}

func (p ConfigParams) IfExists(key string, then func(value string)) {
	v, ok := p[key]
	if ok {
		then(v)
	}
}

func (p ConfigParams) IfNotEmpty(key string, then func(value string)) {
	v, ok := p[key]
	if ok && len(v) > 0 {
		then(v)
	}
}

func (p ConfigParams) IfNotEmptyWithError(key string, then func(value string) error) error {
	v, ok := p[key]
	if ok && len(v) > 0 {
		return then(v)
	}
	return nil
}
