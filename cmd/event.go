package main

import (
	"flag"
	"fmt"
	"foot/event/cmd/app"
	"foot/event/pkg/flags"
	"foot/event/pkg/signal"
	"github.com/spf13/viper"
	"os"
	"time"
)

var(
	argSources     flags.Uri
)

//初始化客户端
func init() {

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("D:/go_workspace/src/foot/event/etc") //改为环境变量
	//viper.AddConfigPath("/root/go/src/foot/event/etc") //改为环境变量
	viper.AddConfigPath(".")
	viper.SetDefault("kubeconfig", "")
	viper.SetDefault("sinks", "mysql")
	viper.SetDefault("interval", time.Second*5)
	flag.Var(&argSources, "source", "source(s) to read events from")
	flag.Parse()
}

func main() {
	command := app.NewEventCommand(signal.SetupSignalHandler())
	//flag.Parse()
	fmt.Println("key:",argSources.Key,"value",argSources.Val)

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
