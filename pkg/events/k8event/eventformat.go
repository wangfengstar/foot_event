package events

import (
	"encoding/json"
	"fmt"
	"foot/event/pkg/options"
	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"reflect"
	"time"
)

func toNodeEvent(event watch.Event) *options.NodeEvent {
	bean := &options.NodeEvent{
		ID:          uuid.New().String(),
		EventType:   string(event.Type),
		Metadata:    toString(event.Object),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	objType := reflect.TypeOf(event.DeepCopy().Object).String()

	fmt.Println(objType)

	if objType == "*v1.Node" {
		if node, ok := event.DeepCopy().Object.(*v1.Node); ok {
			bean.ObjType = "node" //node.Kind
			bean.ObjName = node.Status.Addresses[0].Address
		}
	}

	return bean
}

//wangfeng deployment
func toDeployEvent(event watch.Event) *options.DeployEvent {
	bean := &options.DeployEvent{
		ID:          uuid.New().String(),
		EventType:   string(event.Type),
		Metadata:    toString(event.Object),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	objType := reflect.TypeOf(event.DeepCopy().Object).String()

	res, err := simplejson.NewJson([]byte(bean.Metadata))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	name, err := res.Get("metadata").Get("name").String()
	namespace, err := res.Get("metadata").Get("namespace").String()

	//fmt.Println(namespace)
	//fmt.Println(toString(event.Object))

	if objType == "*v1.Deployment" {
			bean.ObjType = "deployment"
			bean.ObjName = name
			bean.NameSpace = namespace
	}

	return bean
}

//wangfeng pod
func toPodEvent(event watch.Event) *options.PodEvent {
	bean := &options.PodEvent{
		ID:        uuid.New().String(),
		EventType: string(event.Type),
		Metadata:		toString(event.Object),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	objType := reflect.TypeOf(event.DeepCopy().Object).String()

	//fmt.Println(objType)
	//fmt.Println(bean)
	res, err := simplejson.NewJson([]byte(bean.Metadata))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	name, err := res.Get("metadata").Get("name").String()
	namespace, err := res.Get("metadata").Get("namespace").String()

	if objType == "*v1.Pod" {
			bean.ObjType = "pod" //node.Kind
			bean.ObjName = name
			bean.NameSpace = namespace
	}
	return bean
}

func toString(event interface{}) string {
	data, err := json.Marshal(event)
	if err != nil {
		return ""
	} else {
		return string(data)
	}
}
