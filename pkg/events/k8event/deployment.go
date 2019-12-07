package events

import (
	"foot_event/pkg/options"
	"foot_event/pkg/sinks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func StartDeploymentEvent(ctx options.Context) {
	sink := sinks.ManufactureSink()

	//watcher,err := ctx.Client.ExtensionsV1beta1().Deployments("").Watch(metav1.ListOptions{})
	watcher, err := ctx.Client.AppsV1().Deployments("").Watch(metav1.ListOptions{})
	//watcher, err := ctx.Client.CoreV1().Events("").Watch(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for {
		select {
		case event, chanOk := <-watcher.ResultChan():
			if chanOk {
				appEvent := toDeployEvent(event)
				appEvent.ClusterId = ctx.ClusterId
				sink.Update(appEvent)
			}
		default:
			//fmt.Println("do nothing")
		}
		time.Sleep(ctx.Interval)
	}
	<-ctx.Stop
}

//func
