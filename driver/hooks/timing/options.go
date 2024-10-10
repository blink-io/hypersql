package timing

type Option func(*Hook)

func Logf(logf func(string, ...any)) Option {
	return func(h *Hook) {
		h.logf = logf
	}
}
