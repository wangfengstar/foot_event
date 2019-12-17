package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"foot_event/pkg/options"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"os"
	"strings"
	"sync"
	"time"
)

var once sync.Once

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func NewEventCommand(stopCh <-chan struct{}) *cobra.Command {
	//cleanFlagSet := pflag.NewFlagSet(componentEvent, pflag.ContinueOnError)

	/*ops, err := NewEventOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}
	*/

	cmd := &cobra.Command{
		Use:                componentEvent,
		Long:               `Event  is responsible for collecting data.include resource of node and event of cluster.`,
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			eventContext, err := createEventContext(stopCh)
			if err == nil {
				run(eventContext)
			} else {
				//klog.Fatalf("error building event context: %v", err)
			}
		},
	}

	return cmd
}

func run(ctx options.Context) {
	//go func() {
	//	glog.Info("Starting prometheus metrics.")
	//	http.Handle("/metrics", promhttp.Handler())
	//	glog.Warning(http.ListenAndServe(*addr, nil))
	//}()

	glog.Info("Starting all event metrics.")

	events := NewEventInitializers()

	for  event := range events {
		go func(event string) {
			fmt.Printf("event : %v \n", event)

			events[event](ctx)
		}(event)
	}

	s := <-ctx.Stop
	fmt.Printf("event has stoped due to received signal : %v", s)
}

func createEventContext(stop <-chan struct{}) (options.Context, error) {
	var config *rest.Config
	var err error

	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	viper.BindEnv("kubeconfig")
	if forceCfg := os.Getenv("EVENTROUTER_CONFIG"); forceCfg != "" {
		viper.SetConfigFile(forceCfg)
	}

	kubeconfig := viper.GetString("kubeconfig")
	kubeconfig = strings.Replace(kubeconfig, " ", "", -1)
	kubeconfig = strings.Replace(kubeconfig, " ", "", -1)
	if len(kubeconfig) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	/*	cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}})*/
	loader := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	oconfig, _ := loader.Load()
	/*	fmt.Println(toString(ss.Clusters))
		fmt.Println(ss.Contexts)
		fmt.Println(ss.AuthInfos)
		fmt.Println(ss.Preferences)*/
	if err != nil {
		klog.Error("init kubernetes client config error :", err.Error())
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		klog.Error("init kubernetes client error :", err.Error())
		panic(err.Error())
	}
	interval := viper.GetDuration("interval") * time.Second
	cluster := viper.GetString("cluster")
	devs := viper.GetStringSlice("devs")

	sharedInformers := informers.NewSharedInformerFactory(clientset, interval)
	ctx := options.Context{
		ClusterId:       cluster,
		InformerFactory: sharedInformers,
		Devs:            devs,
		Client:          clientset,
		Kubeconfig:      config,
		Stop:            stop,
		Interval:        interval,
		Config:          oconfig,
	}
	return ctx, nil
}

func toString(event interface{}) string {
	data, err := json.Marshal(event)
	if err != nil {
		return ""
	} else {
		return string(data)
	}
}
