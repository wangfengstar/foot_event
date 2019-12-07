package sinks

import (
	"fmt"
	"foot_event/pkg/options"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/watch"
	"reflect"
	"sync"
	"time"
)

type MysqlSink struct {
	config MysqlConfig
	client *xorm.Engine
	sync.RWMutex
	dbExists bool
}

type MysqlConfig struct {
	ClusterName     string
	User            string
	Password        string
	Host            string
	Port            int
	DbName          string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	Secure          bool
}

// Returns a thread-safe implementation of EventSinkInterface for InfluxDB.
func NewMysqlSink(cfg MysqlConfig) (EventSinkInterface, error) {
	client, err := newMysqlClient(cfg)
	if err != nil {
		return nil, err
	}

	return &MysqlSink{
		config:   cfg,
		client:   client,
		dbExists: false,
	}, nil
}

func newMysqlClient(c MysqlConfig) (*xorm.Engine, error) {
	//root:password@tcp(localhost:3306)/database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", c.User, c.Password, c.Host, c.Port, c.DbName)
	engine, err := xorm.NewEngine("mysql", dsn)
	//defer engine.
	if err != nil {
		fmt.Printf("create mysql engine,err:%v\n", err)
		return nil, err
	}
	engine.SetMaxOpenConns(c.MaxOpenConns)
	engine.SetMaxIdleConns(c.MaxIdleConns)

	//logWriter, err := os.Create("F:/sql.log")
	//if err != nil {
	//	log.Fatalf("Fail to create xorm system logger: %v\n", err)
	//}
	//logger := xorm.NewSimpleLogger(logWriter)
	//logger.ShowSQL(true)
	//engine.SetLogger(logger)
	return engine, nil
}

func (sink *MysqlSink) Insert(obj interface{}) {
	sink.Lock()
	defer sink.Unlock()
	sink.insert(obj)
}

//一个集群中的Internel IP是唯一的
var nodeMap map[string]string = make(map[string]string)
var imageMap map[string]string = make(map[string]string)

func (sink *MysqlSink) Update(obj interface{}) {
	sink.Lock()
	defer sink.Unlock()
	objType := reflect.TypeOf(obj).String()
	if objType == "[]*options.Node" {
		if nodes, ok := obj.([]*options.Node); ok {
			for _, node := range nodes {
				if _, ok := nodeMap[node.IP]; ok {
					node.UpdatedTime = time.Now()
					sink.client.Update(node, &options.Node{ClusterId: node.ClusterId, IP: node.IP})
				} else {
					sink.insert(node)
					nodeMap[node.IP] = "1"
				}

			}
		}
	} else if objType == "[]*options.Image" {
		if images, ok := obj.([]*options.Image); ok {
			for _, image := range images {
				if _, ok := imageMap[image.Name]; ok {
					image.UpdatedTime = time.Now()
					sink.client.Update(image, &options.Image{ClusterId: image.ClusterId, Host: image.Host, Name: image.Name})
				} else {
					sink.insert(image)
					imageMap[image.Name] = "1"
				}

			}
		}
	} else if objType == "*options.Cluster" {
		if Cluster, ok := obj.(*options.Cluster); ok {
			sink.client.Update(Cluster)
		}
	} else {
		sink.insert(obj)
	}
}

// UpdateEvents implements the EventSinkInterface
func (sink *MysqlSink) ExportNodeEvents(event watch.Event) {

}

func (sink *MysqlSink) insert(node interface{}) error {
	if sink.client == nil {
		client, err := newMysqlClient(sink.config)
		if err != nil {
			return err
		}
		sink.client = client
	}
	err := sink.client.Ping()
	if err != nil {
		glog.Infof("can not ping database server ")
		return err
	}
	affected, existErr := sink.client.Insert(node)
	if existErr != nil {
		glog.Infof("export data to database server,err :%v ", existErr)
		return existErr
	}
	glog.Infof("sync node size : %v", affected)
	return nil
}
