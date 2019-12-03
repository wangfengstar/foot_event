package events

//wangfeng

import (
	"foot/event/pkg/options"
	"foot/event/pkg/sinks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func StartPodEvent(ctx options.Context) {
	sink := sinks.ManufactureSink()
	//watcher, err := ctx.Client.CoreV1().Pods("").Watch(metav1.ListOptions{})
	watcher, err := ctx.Client.CoreV1().Pods("").Watch(metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for {
		select {
		case event, chanOk := <-watcher.ResultChan():
			if chanOk {
				nodeEvent := toPodEvent(event)
				nodeEvent.ClusterId = ctx.ClusterId
				sink.Update(nodeEvent)
			}
		default:
			//fmt.Println("do nothing")
		}
		time.Sleep(ctx.Interval)
	}
	<-ctx.Stop
}
