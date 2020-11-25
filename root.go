package rap

import (
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogger *zap.Logger

	rootLogger *Logger

	defaultConfig = zapcore.EncoderConfig{
		MessageKey:   "msg",                       //结构化（json）输出：msg的key
		LevelKey:     "level",                     //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",                        //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                      //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
)

func init() {
	defaultCore := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(defaultConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zap.DebugLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
	)

	zapLogger = zap.New(defaultCore, zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(zap.WarnLevel))

	rootLogger = &Logger{zapLogger}
}

// 初始化日志 logger
func InitLog(logPath, errPath string, logLevel zapcore.Level, writer func(string) io.Writer) {
	config := zapcore.EncoderConfig{
		MessageKey:   "msg",                       //结构化（json）输出：msg的key
		LevelKey:     "level",                     //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",                        //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                      //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
		}, //输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
	}
	//自定义日志级别：自定义Info级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})

	//自定义日志级别：自定义Warn级别
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	// 获取io.Writer的实现
	infoWriter := writer(logPath)
	warnWriter := writer(errPath)

	// 实现多个输出
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(infoWriter), infoLevel),                               //将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
		zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(warnWriter), warnLevel),                               //warn及以上写入errPath
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
	)

	zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(zap.WarnLevel))
	rootLogger = &Logger{zapLogger}
}

func GetLumberjackWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,    //最大M数，超过则切割
		MaxBackups: 5,     //最大文件保留数，超过就删除最老的日志文件
		MaxAge:     30,    //保存30天
		Compress:   false, //是否压缩
	}
}

func Sync() {
	zapLogger.Sync()
}

func New(name string) *Logger {
	return rootLogger.Named(name)
}
func Debug(msg string, ctx ...interface{}) {
	rootLogger.Debug(msg, ctx...)
}
func Info(msg string, ctx ...interface{}) {
	rootLogger.Info(msg, ctx...)
}
func Warn(msg string, ctx ...interface{}) {
	rootLogger.Warn(msg, ctx...)
}
func Error(msg string, ctx ...interface{}) {
	rootLogger.Error(msg, ctx...)
}
func DPanic(msg string, ctx ...interface{}) {
	rootLogger.DPanic(msg, ctx...)
}
func Panic(msg string, ctx ...interface{}) {
	rootLogger.Panic(msg, ctx...)
}
func Fatal(msg string, ctx ...interface{}) {
	rootLogger.Fatal(msg, ctx...)
}
