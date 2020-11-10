package rap

import (
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"testing"
)

func TestRap(t *testing.T) {
	var writer = func(filename string) io.Writer {
		return &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    10,    //最大M数，超过则切割
			MaxBackups: 5,     //最大文件保留数，超过就删除最老的日志文件
			MaxAge:     30,    //保存30天
			Compress:   false, //是否压缩
		}
	}

	InitLog("./info.log", "./error.log", zap.InfoLevel, writer)
	defer Sync()
	l2 := New("main")
	l2.Debug("hello, world", "name", "wang", "age", 33)
	l2.Info("hello, world", "name", "wang", "age", 33)
	l2.Warn("hello, world", "name", "wang", "age", 33)
	l2.Error("hello, world", "name", "wang", "age", 33)
	l2.Debug("hello, world", "name", "wang", "age", 33)
	l2.DPanic("hello, world", "name", "wang", "age", 33)
}

