package app

import (
	evtk8 "foot_event/pkg/events/k8event"
)

//var wg sync.WaitGroup

func NewEventInitializers() map[string]InitFunc{
	events := map[string]InitFunc{}
	//fmt.Println("do nothing")
	//events["node"] = evtk8.StartNodeEvent
	events["scheduler"] = evtk8.StartSchedulerEvent
	events["events"] = evtk8.StartEvents
	events["pod"] = evtk8.StartPodEvent
	//events["NS"] = evt.TestNodeEvent
	//events["nodeStat"] = evtk8.InitNodeStats
	events["deployment"] = evtk8.StartDeploymentEvent
	//events["mesosStat"] = evtmesos.MesosEvents
	//wg.Add(20)

	return events
}
