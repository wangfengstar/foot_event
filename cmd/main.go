package main

import (
	"flag"
	"fmt"
	"foot/event/cmd/app"
	"foot/event/pkg/signal"
	"github.com/spf13/viper"
	"os"
	"time"
)

//初始化客户端
func init() {
	flag.Parse()
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("D:/go_workspace/src/foot/event/etc") //改为环境变量
	//viper.AddConfigPath("/root/go/src/foot/event/etc") //改为环境变量
	viper.AddConfigPath(".")
	viper.SetDefault("kubeconfig", "")
	viper.SetDefault("sinks", "mysql")
	viper.SetDefault("interval", time.Second*5)
}

func main() {
	command := app.NewEventCommand(signal.SetupSignalHandler())

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
