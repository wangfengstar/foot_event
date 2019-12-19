package events

//wangfeng

import (
	"fmt"
	"foot_event/pkg/loggs"
	"foot_event/pkg/options"
	"foot_event/pkg/sinks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func StartSchedulerEvent(ctx options.Context) {
	sink := sinks.ManufactureSink()

	watcher, err := ctx.Client.AutoscalingV1().HorizontalPodAutoscalers("").Watch(metav1.ListOptions{})
	if err != nil {
		loggs.Log.Error(fmt.Sprint("error HorizontalPodAutoscalers watcher: ", err))
		panic(err.Error())
	}

	for {
		select {
		case event, chanOk := <-watcher.ResultChan():
			if chanOk {
				nodeEvent := toSchedulerEvent(event)
				nodeEvent.ClusterId = ctx.ClusterId
				sink.Update(nodeEvent)
			}
		default:
			loggs.Log.Info(fmt.Sprint("StartSchedulerEvent do nothing!"))
		}
		time.Sleep(ctx.Interval)
	}
	<-ctx.Stop
}
