package main

import (
	"flag"
	"fmt"
	"foot_event/cmd/app"
	"foot_event/pkg/loggs"
	"foot_event/pkg/signal"
	"github.com/spf13/viper"
	"os"
	"time"
)

//初始化客户端
func init() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath( os.Getenv("WORK_HOME")+"/etc") //改为环境变量
	viper.AddConfigPath(".")
	viper.SetDefault("kubeconfig", "")
	viper.SetDefault("sinks", "mysql")
	viper.SetDefault("interval", time.Second*5)
	//日志初始化
	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	loggs.InitLogger(os.Getenv("WORK_HOME")+viper.GetString("log_dir")+viper.GetString("log_file"),viper.GetString("log_level"))

	flag.Parse()
}

func main() {
	command := app.NewEventCommand(signal.SetupSignalHandler())

	loggs.Log.Info(fmt.Sprint("foot_event start..... "))
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		loggs.Log.Error(fmt.Sprint("foot_event exit: ", err))
		os.Exit(1)
	}
}
