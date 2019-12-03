package events

import (
	"fmt"
	"foot/event/pkg/options"
	"foot/event/pkg/sinks"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/disk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"sync"
	"time"
)

var once sync.Once

const (
	SHA    = "@sha256:"
	DOMAIN = "ai.foot.com/domain" //"beta.kubernetes.io/arch"
)

func StartNodeEvent(ctx options.Context) {
	sink := sinks.ManufactureSink()

	once.Do(func() {
		fmt.Println("------------")

		/*info := &options.Cluster{
			ID:			ctx.ClusterId,
			ClusterURL: ctx.Kubeconfig.Host,
			CertData:	string(ctx.Kubeconfig.TLSClientConfig.CertData),
			KeyData:    string(ctx.Kubeconfig.TLSClientConfig.KeyData),
			CaData:     string(ctx.Kubeconfig.TLSClientConfig.CAData),
		}*/

		var contexts []map[string]interface{}
		for name, context := range ctx.Config.Contexts {
			cmap := make(map[string]interface{})
			cmap["name"] = name
			cmap["context"] = context
			contexts = append(contexts, cmap)
		}
		var clusters []map[string]interface{}
		for name, cluster := range ctx.Config.Clusters {
			cmap := make(map[string]interface{})
			cmap["name"] = name
			cmap["cluster"] = cluster
			clusters = append(clusters, cmap)
		}
		var users []map[string]interface{}
		for name, user := range ctx.Config.AuthInfos {
			cmap := make(map[string]interface{})
			cmap["name"] = name
			cmap["user"] = user
			users = append(users, cmap)
		}
		//fmt.Println(toString(ctx.Config.Contexts))
		info := &options.Cluster{
			ID:         ctx.ClusterId,
			ClusterURL: ctx.Kubeconfig.Host,
			Contexts:   toString(contexts),
			Clusters:   toString(clusters),
			Users:      toString(users),
		}
		sink.Update(info)
	})

	watcher, err := ctx.Client.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for {
		select {
		case event, chanOk := <-watcher.ResultChan():
			if chanOk {
				//fmt.Println(toString(event.Object))
				nodeEvent := toNodeEvent(event)
				nodeEvent.ClusterId = ctx.ClusterId
				//fmt.Println(nodeEvent.ClusterId)
				sink.Update(nodeEvent)
			}
		default:
			fmt.Println("do nothing")
		}
		time.Sleep(ctx.Interval)
	}
	<-ctx.Stop
}

func InitNodeStats(ctx options.Context) {
	sink := sinks.ManufactureSink()
	for {
		list, err := ctx.Client.CoreV1().Nodes().List(metav1.ListOptions{})
		nodes := make([]*options.Node, 0)
		if err == nil {
			for _, item := range list.Items {
				id := uuid.New().String()

				ip := item.Status.Addresses[0].Address
				node := &options.Node{
					ID:             id, //uuid.New().String(),
					ClusterId:      ctx.ClusterId,
					Host:           item.Name,
					SystemVersion:  item.Status.NodeInfo.KubeletVersion,
					Kernel:         item.Status.NodeInfo.KernelVersion,
					DockerVersion:  item.Status.NodeInfo.ContainerRuntimeVersion,
					KubeletVersion: item.Status.NodeInfo.KubeletVersion,
					Status:         toString(item.Status.Conditions),
					Domain:         "default",
					IP:             ip,
					ExtIP:          "",
					CPU:            item.Status.Allocatable.Cpu().Size(),
					Mem:            item.Status.Allocatable.Memory().Size(),
					Disk:           diskStats(ctx.Devs, ip),
					GPU:            0,
					Labels:         toString(item.Labels),
					Annotations:    toString(item.Annotations),
					PodCidr:        item.Spec.PodCIDR,
					Pods:           item.Status.Allocatable.Pods().Size(),
					CreatedTime:    time.Now(),
					UpdatedTime:    time.Now(),
				}
				if value, ok := item.Labels[DOMAIN]; ok {
					node.Domain = value
				}
				nodes = append(nodes, node)
				images := make([]*options.Image, 0)
				for _, img := range item.Status.Images {
					for _, name := range img.Names {
						if !strings.Contains(name, SHA) {
							id := uuid.New().String()
							image := &options.Image{
								ID:          id,
								ClusterId:   ctx.ClusterId,
								Host:        item.Name,
								Name:        name,
								Size:        img.SizeBytes,
								CreatedTime: time.Now(),
								UpdatedTime: time.Now(),
							}
							images = append(images, image)
						}

					}
				}
				sink.Update(images)
				fmt.Printf(" init node [%v] meta info to sinks\n", item.Name)
			}
			sink.Update(nodes)
		} else {
			fmt.Errorf("init node meta info err : %v", err)
		}
		time.Sleep(ctx.Interval)
	}
}

//TODO：数据占用太大空间时，可转换成必要的有效数据入库
func diskStats(disks []string, ip string) string {
	stats := []*disk.UsageStat{}
	for _, path := range disks {
		if stat, err := disk.Usage(path); err == nil {
			stats = append(stats, stat)
		} else {
			fmt.Printf("warn: %v not exist on host [%v] \n", path, ip)
		}
	}

	return toString(stats)
}
