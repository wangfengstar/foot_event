/*
Copyright 2017 Heptio Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sinks

import (
	"errors"

	"k8s.io/apimachinery/pkg/watch"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	_ "github.com/go-sql-driver/mysql"
)

// EventSinkInterface is the interface used to shunt events
type EventSinkInterface interface {
	Insert(obj interface{})
	Update(obj interface{})
	ExportNodeEvents(event watch.Event)
}

// ManufactureSink will manufacture a sinks according to viper configs
// TODO: Determine if it should return an array of sinks
func ManufactureSink() (e EventSinkInterface) {
	s := viper.GetString("sinks")
	glog.Infof("Sink is [%v]", s)
	switch s {
	case "glog":
		e = NewGlogSink()
	case "mysql":
		host := viper.GetString("host")
		if host == "" {
			panic("mysql sinks specified but host not specified")
		}

		username := viper.GetString("username")
		if username == "" {
			panic("mysql sinks specified but username not specified")
		}

		password := viper.GetString("password")
		if password == "" {
			panic("mysql sinks specified but password not specified")
		}

		port := viper.GetInt("port")
		if password == "" {
			panic("mysql sinks specified but port not specified")
		}
		dbName := viper.GetString("dbName")
		if password == "" {
			panic("mysql sinks specified but dbName not specified")
		}

		cluterName := viper.GetString("cluterName")
		if cluterName == "" {
			cluterName = "default"
			glog.V(4).Infof("mysql sinks specified but cluster not specified,default: %s","default")
		}

		connMaxLifetime := viper.GetDuration("connMaxLifetime")
		if connMaxLifetime == 0 {
			connMaxLifetime = 20
			glog.V(4).Infof("mysql sinks specified but cluster not specified,default: %s","default")
		}

		maxOpenConns := viper.GetInt("maxOpenConns")
		if maxOpenConns == 0 {
			maxOpenConns = 20
			glog.V(4).Infof("mysql sinks specified but cluster not specified,default: %s","default")
		}

		maxIdleConns := viper.GetInt("maxIdleConns")
		if maxIdleConns == 0 {
			maxIdleConns = 20
			glog.V(4).Infof("mysql sinks specified but cluster not specified,default: %s","default")
		}

		secure := viper.GetBool("secure")


		cfg := MysqlConfig{
			User:                  username,
			Password:              password,
			Secure:                secure,
			Host:                  host,
			Port:                  port,
			DbName:                dbName,
			ConnMaxLifetime:	   connMaxLifetime,
			MaxOpenConns:		   maxOpenConns,
			MaxIdleConns:		   maxIdleConns,
			ClusterName:           cluterName,
		}

		mysql, err := NewMysqlSink(cfg)
		if err != nil {
			panic(err.Error())
		}
		return mysql
	default:
		err := errors.New("Invalid Sink Specified")
		panic(err.Error())
	}
	return e
}
