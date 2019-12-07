package app

import (
	"foot_event/pkg/options"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	componentEvent = "event" //组件名称
)

type InitFunc func(ctx options.Context)

type EventOption struct {
	KubeConfig string
	Interval   time.Duration
	Sink       string
}

func NewEventOptions() (*EventOption, error) {

	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	viper.BindEnv("kubeconfig")
	if forceCfg := os.Getenv("EVENTROUTER_CONFIG"); forceCfg != "" {
		viper.SetConfigFile(forceCfg)
	}
	eventOps := EventOption{
		KubeConfig: viper.GetString("kubeconfig"),
		Interval:   viper.GetDuration("interval"),
		Sink:       viper.GetString("sinks"),
	}

	return &eventOps, nil
}
