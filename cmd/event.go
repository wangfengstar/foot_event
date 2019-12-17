package main

import (
	"flag"
	"fmt"
	"foot_event/cmd/app"
	"foot_event/pkg/flags"
	"go.uber.org/zap/zapcore"
	"github.com/natefinch/lumberjack"
	"foot_event/pkg/signal"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"time"
)

var(
	argSources     flags.Uri
)

func initLogger(logpath string, loglevel string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径
		MaxSize:    128,     // megabytes
		MaxBackups: 30,      // 最多保留300个备份
		MaxAge:     7,       // days
		Compress:   true,    // 是否压缩 disabled by default
	}
	w := zapcore.AddSync(&hook)

	var level zapcore.Level
	switch loglevel {
		case "debug":
			level = zap.DebugLevel
		case "info":
			level = zap.InfoLevel
		case "error":
			level = zap.ErrorLevel
		default:
			level = zap.InfoLevel
		}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,)

	logger := zap.New(core)
	logger.Info("DefaultLogger init success")

	return logger
}

//初始化客户端
func init() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath( os.Getenv("WORK_HOME")+"/foot_event/etc") //改为环境变量
	viper.AddConfigPath(".")
	viper.SetDefault("kubeconfig", "")
	viper.SetDefault("sinks", "mysql")
	viper.SetDefault("interval", time.Second*5)

	flag.Var(&argSources, "source", "source(s) to read events from")
	flag.Parse()
}

func main() {
	logger := initLogger("./event.log", "info")
	command := app.NewEventCommand(signal.SetupSignalHandler())

	logger.Info(fmt.Sprint("event start..... "))
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)

		os.Exit(1)
	}

}
