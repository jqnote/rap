package rap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"reflect"
)

type Logger struct {
	l *zap.Logger
}

func (log Logger) Named(name string) *Logger {
	l := log.l.Named(name)
	return &Logger{l}
}

func (log *Logger) Debug(msg string, ctx ...interface{}) {
	log.write(zap.DebugLevel, msg, ctx)
}

func (log *Logger) Info(msg string, ctx ...interface{}) {
	log.write(zap.InfoLevel, msg, ctx)
}

func (log *Logger) Warn(msg string, ctx ...interface{}) {
	log.write(zap.WarnLevel, msg, ctx)
}

func (log *Logger) Error(msg string, ctx ...interface{}) {
	log.write(zap.ErrorLevel, msg, ctx)
}

func (log *Logger) DPanic(msg string, ctx ...interface{}) {
	log.write(zap.DPanicLevel, msg, ctx)
}

func (log *Logger) Panic(msg string, ctx ...interface{}) {
	log.write(zap.PanicLevel, msg, ctx)
}

func (log *Logger) Fatal(msg string, ctx ...interface{}) {
	log.write(zap.FatalLevel, msg, ctx)
}

func (log *Logger) write(lvl zapcore.Level, msg string, args []interface{}) {
	if ce := log.l.Check(lvl, msg); ce != nil {
		var fields []zap.Field
		size := len(args)
		for i := 0; i < size; i = i + 2 {
			key := args[i].(string)
			if size == i+1 {
				fields = append(fields, zap.String(key, ""))
			} else {
				val := args[i+1]
				typ := reflect.TypeOf(val)
				switch typ.Kind() {
				case reflect.String:
					fields = append(fields, zap.String(key, args[i+1].(string)))
				case reflect.Bool:
					fields = append(fields, zap.Bool(key, val.(bool)))
				case reflect.Int:
					fields = append(fields, zap.Int(key, val.(int)))
				case reflect.Int8:
					fields = append(fields, zap.Int8(key, val.(int8)))
				case reflect.Int16:
					fields = append(fields, zap.Int16(key, val.(int16)))
				case reflect.Int32:
					fields = append(fields, zap.Int32(key, val.(int32)))
				case reflect.Int64:
					fields = append(fields, zap.Int64(key, val.(int64)))
				case reflect.Uint:
					fields = append(fields, zap.Uint(key, val.(uint)))
				case reflect.Uint8:
					fields = append(fields, zap.Uint8(key, val.(uint8)))
				case reflect.Uint16:
					fields = append(fields, zap.Uint16(key, val.(uint16)))
				case reflect.Uint32:
					fields = append(fields, zap.Uint32(key, val.(uint32)))
				case reflect.Uint64:
					fields = append(fields, zap.Uint64(key, val.(uint64)))
				case reflect.Uintptr:
					fields = append(fields, zap.Uintptr(key, val.(uintptr)))
				case reflect.Float32:
					fields = append(fields, zap.Float32(key, val.(float32)))
				case reflect.Float64:
					fields = append(fields, zap.Float64(key, val.(float64)))
				case reflect.Complex64:
					fields = append(fields, zap.Complex64(key, val.(complex64)))
				case reflect.Complex128:
					fields = append(fields, zap.Complex128(key, val.(complex128)))
				default:
					fields = append(fields, zap.Any(key, val))
				}
			}
		}
		ce.Write(fields...)
	}
}

