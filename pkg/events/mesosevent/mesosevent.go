package events

import (
	"fmt"
	"foot_event/pkg/options"
	"io/ioutil"
	"net/http"
	"strings"
)

type LoginInfo struct {
	Username string
	Password string
	LoginYzm string
}

func MesosEvents(ctx options.Context) {

	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET",
			"http://aihps:ai2016@10.1.245.19:8080/v2/groups",
			strings.NewReader(""))
		if err != nil {
			fmt.Println("error...")
			return
		}
		req.Header.Set("accept", "text/event-stream")
		resp, err := client.Do(req)
		if err != nil {
			return
		}

		if resp.StatusCode == 200 {
			fmt.Println("true...")
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error...")
			return
		}
		fmt.Println(string(body))
	}

}

//func MesosEventsStream(ctx Context) {
//	//sink := sinks.ManufactureSink()
//	client := &http.Client{}
//	req, err := http.NewRequest("GET",
//		"http://aihps:ai2016@10.1.245.19:8080/v2/events",
//		strings.NewReader(""))
//	if err != nil {
//		fmt.Println("error...")
//		return
//	}
//	req.Header.Set("accept", "text/event-stream")
//	resp, err := client.Do(req)
//	defer resp.Body.Close()
//
//	for {
//		select {
//			case body, err := ioutil.ReadAll(resp.Body)
//			if err != nil {
//				fmt.Println("error...")
//				return
//			}
//			fmt.Println(string(body))
//		}
//	}
//}
