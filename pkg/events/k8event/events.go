package events

import (
	"fmt"
	"foot_event/pkg/loggs"
	"foot_event/pkg/options"
	"foot_event/pkg/sinks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

//wangfeng
func StartEvents(ctx options.Context) {
	sink := sinks.ManufactureSink()
	//watcher, err := ctx.Client.CoreV1().Events("").Watch(metav1.ListOptions{})
	watcher, err := ctx.Client.CoreV1().Events("").Watch(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for {
		select {
		case event, chanOk := <-watcher.ResultChan():
			if chanOk {
				nodeEvent := toEvents(event)
				nodeEvent.ClusterId = ctx.ClusterId
				sink.Update(nodeEvent)
			}
		default:
			loggs.Log.Info(fmt.Sprint("StartEvents do nothing!"))
		}
		time.Sleep(ctx.Interval)
	}
	<-ctx.Stop
}

