package hypersql

var commonErrorHandlers = make(map[error]func(error) *Error)

func RegisterCommonErrorHandler(e error, f func(error) *Error) {
	commonErrorHandlers[e] = f
}

func handleCommonError(e error) (func(error) *Error, bool) {
	fn, ok := commonErrorHandlers[e]
	return fn, ok
}
