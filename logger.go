package hypersql

type Logger func(format string, args ...any)

func NoopLogger(format string, args ...any) {}

//func doLoggerFunc(l Logger) func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
//	return func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
//		var sb strings.Builder
//		for k, v := range data {
//			sb.WriteString(fmt.Sprintf(" %s=%v", k, v))
//		}
//		l("msg: %s, data:%s", msg, sb.String())
//	}
//}
