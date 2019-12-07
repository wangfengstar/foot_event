package events

//wangfeng

import (
	"fmt"
	"foot_event/pkg/options"
	"foot_event/pkg/sinks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func StartSchedulerEvent(ctx options.Context) {
	sink := sinks.ManufactureSink()
	watcher, err := ctx.Client.AppsV1().Deployments("default").Watch(metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for {
		select {
		case event, chanOk := <-watcher.ResultChan():
			if chanOk {
				nodeEvent := toNodeEvent(event)
				nodeEvent.ClusterId = ctx.ClusterId
				sink.Update(nodeEvent)
			}
		default:
			fmt.Println("do nothing")
		}
		time.Sleep(ctx.Interval)
	}
	<-ctx.Stop
}
