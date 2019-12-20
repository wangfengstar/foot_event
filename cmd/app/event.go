package app

import (
	"fmt"
	evtk8 "foot_event/pkg/events/k8event"
	"foot_event/pkg/loggs"
)


func NewEventInitializers() map[string]InitFunc{
	events := map[string]InitFunc{}
	loggs.Log.Info(fmt.Sprint("foot_event load event type start ..... "))
	//events["node"] = evtk8.StartNodeEvent
	events["scheduler"] = evtk8.StartSchedulerEvent
	events["events"] = evtk8.StartEvents
	events["pod"] = evtk8.StartPodEvent
	//events["NS"] = evt.TestNodeEvent
	//events["nodeStat"] = evtk8.InitNodeStats
	events["deployment"] = evtk8.StartDeploymentEvent
	//events["mesosStat"] = evtmesos.MesosEvents
	loggs.Log.Info(fmt.Sprint("foot_event load event type end ..... "))

	return events
}
