package main

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"encoding/json"
)
func main(){
	stats := []*disk.UsageStat{}
	disks := []string{"D:/"}
	for _,path := range disks{
		if stat, err := disk.Usage(path);err == nil {
			stats = append(stats, stat)
		}else{
			fmt.Printf("warn: %v not exist on host [%v] \n",path,"10.1.235.89")
		}
	}

	fmt.Println(toString(stats))
}

func toString(event interface{})string{
	data,err := json.Marshal(event)
	if err != nil {
		return ""
	}else{
		return string(data)
	}
}
