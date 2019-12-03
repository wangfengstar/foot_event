package options

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"time"
)

type Context struct {
	ClusterId       string
	InformerFactory informers.SharedInformerFactory
	Client          *kubernetes.Clientset
	Kubeconfig      *rest.Config
	Interval        time.Duration
	Devs            []string
	Config          *clientcmdapi.Config
	Stop            <-chan struct{}
}

type Node struct {
	ID             string    `xorm:"'ID'"`
	ClusterId      string    `xorm:"'CLUSTER_ID'"`      /** 资源池;集群ID */
	Host           string    `xorm:"'HOST'"`            /** 主机名称 */
	SystemVersion  string    `xorm:"'SYSTEM_VERSION'"`  /** 系统版本 */
	Kernel         string    `xorm:"'KERNEL'"`          /** 内核信息 */
	DockerVersion  string    `xorm:"'DOCKER_VERSION'"`  /** docker版本 */
	KubeletVersion string    `xorm:"'KUBELET_VERSION'"` /** kubelet版本 */
	Status         string    `xorm:"'STATUS'"`          /** 主机状态 */
	Domain         string    `xorm:"'DOMAIN'"`          /** 节点所属域;节点所属域 */
	IP             string    `xorm:"'IP'"`              /** 默认网卡（管理） */
	ExtIP          string    `xorm:"'EXT_IP'"`          /** 扩展网卡 */
	CPU            int       `xorm:"'CPU'"`             /** CPU;可分配 */
	Mem            int       `xorm:"'MEM'"`             /** MEM;可分配 */
	Disk           string    `xorm:"'DISK'"`            /** DISK;可分配 */
	GPU            int       `xorm:"'GPU'"`             /** GPU */
	Labels         string    `xorm:"'LABELS'"`          /** 标签;可分配 */
	Annotations    string    `xorm:"'ANNOTATIONS'"`     /** 注解;注解 */
	PodCidr        string    `xorm:"'POD_CIDR'"`        /** pod地址;podIP */
	Pods           int       `xorm:"'PODS'"`            /** GPU */
	CreatedTime    time.Time `xorm:"'CREATED_TIME'"`
	UpdatedTime    time.Time `xorm:"'UPDATED_TIME'"`
}
type Image struct {
	ID          string    `xorm:"'ID'"`
	ClusterId   string    `xorm:"'CLUSTER_ID'"` /** 资源池;集群ID */
	Host        string    `xorm:"'HOST'"`       /** 主机 */
	Name        string    `xorm:"'NAME'"`       /** 镜像名称 */
	Size        int64     `xorm:"'SIZE'"`       /** 大小 */
	CreatedTime time.Time `xorm:"'CREATED_TIME'"`
	UpdatedTime time.Time `xorm:"'UPDATED_TIME'"`
}

type NodeEvent struct {
	ID          string    `xorm:"'ID'"`
	ClusterId   string    `xorm:"'CLUSTER_ID'"` /** 资源池;集群ID */
	EventType   string    `xorm:"'EVT_TYPE'"`   /** 事件类型 */
	ObjName     string    `xorm:"'OBJ_NAME'"`   /** 事件对象名称 */
	ObjType     string    `xorm:"'OBJ_TYPE'"`   /** 事件对象 */
	Metadata    string    `xorm:"'META_DATA'"`  /** 元数据 */
	CreatedTime time.Time `xorm:"'CREATED_TIME'"`
	UpdatedTime time.Time `xorm:"'UPDATED_TIME'"`
}

//wangfeng pod
type PodEvent struct {
	ID          string    `xorm:"'ID'"`
	ClusterId   string    `xorm:"'CLUSTER_ID'"` /** 资源池;集群ID */
	EventType   string    `xorm:"'EVT_TYPE'"`   /** 事件类型 */
	NameSpace   string    `xorm:"'NAMESPACE'"`  /** 域名 */
	ObjName     string    `xorm:"'OBJ_NAME'"`   /** 事件对象名称 */
	ObjType     string    `xorm:"'OBJ_TYPE'"`   /** 事件对象 */
	Metadata    string    `xorm:"'META_DATA'"`  /** 元数据 */
	CreatedTime time.Time `xorm:"'CREATED_TIME'"`
	UpdatedTime time.Time `xorm:"'UPDATED_TIME'"`
}

//wangfeng deployment
type DeployEvent struct {
	ID          string    `xorm:"'ID'"`
	ClusterId   string    `xorm:"'CLUSTER_ID'"` /** 资源池;集群ID */
	EventType   string    `xorm:"'EVT_TYPE'"`   /** 事件类型 */
	NameSpace   string    `xorm:"'NAMESPACE'"`  /** 域名 */
	ObjName     string    `xorm:"'OBJ_NAME'"`   /** 事件对象名称 */
	ObjType     string    `xorm:"'OBJ_TYPE'"`   /** 事件对象 */
	Metadata    string    `xorm:"'META_DATA'"`  /** 元数据 */
	CreatedTime time.Time `xorm:"'CREATED_TIME'"`
	UpdatedTime time.Time `xorm:"'UPDATED_TIME'"`
}


type Cluster struct {
	ID          string `xorm:"'ID'"`
	ClusterURL  string `xorm:"'CLUSTER_URL'"` /** 资源池;集群ID */
	Contexts    string `xorm:"'CONTEXTS'"`    /** 事件类型 */
	Clusters    string `xorm:"'CLUSTERS'"`    /** 事件对象名称 */
	Users       string `xorm:"'USERS'"`       /** 事件对象 */
	Preferences string `xorm:"'PREFERENCES'"`
}

func (Cluster) TableName() string {
	return "CM_CLUSTER"
}

func (Node) TableName() string {
	return "CM_HOST"
}

func (Image) TableName() string {
	return "CM_HOST_IMAGE"
}
func (NodeEvent) TableName() string {
	return "CM_EVENT"
}

//wangfeng pod
func (PodEvent) TableName() string {
	return "POD_EVENT"
}

////wangfeng Deploy
func (DeployEvent) TableName() string {
	return "APP_EVENT"
}
