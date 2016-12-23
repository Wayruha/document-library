package main

import (
	"fmt"
	"log"
	"github.com/samuel/go-zookeeper/zk"
	"math/rand"
        "net/http"
        "net/http/httputil"
	"strings"
	"time"
)

var urls []string


func must(err error) {
	if err != nil {
		//panic(err)
		fmt.Printf("%+v From must \n", err)
	}
}

func connect() *zk.Conn {
	zksStr := "172.21.0.2:2181" // Need to replace by service name
	zks := strings.Split(zksStr, ",")
	conn, _, err := zk.Connect(zks, time.Second)
	must(err)
	return conn
}

func monitorGserver(conn *zk.Conn, path string) (chan []string, chan error) {

	servers := make(chan []string)
	errors := make(chan error)
	go func() {
		for {
			children, _, events, err := conn.ChildrenW(path)
			if err != nil {
				errors <- err
				return
			}
			servers <- children
			evt := <-events
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
		}
	}()
	return servers, errors
}

func NewMultipleHostReverseProxy() *httputil.ReverseProxy {
        director := func(req *http.Request) {

		if (req.URL.Path == "/library") {
			fmt.Println("This is for gserver")
			target := urls[rand.Int()%len(urls)]
			req.URL.Scheme = "http"
                	req.URL.Host = target
                	req.URL.Path ="/"
	
		}else {
        	   fmt.Println("This is for nginx")
			req.URL.Scheme = "http"
                	req.URL.Host = "172.21.0.3:80" // Need to replace by service name
                	//req.URL.Path = "/"
    		}
               
        }
        return &httputil.ReverseProxy{Director: director}
}

func main() {
	conn := connect()
	defer conn.Close()

	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)

	exists, stat, err := conn.Exists("/grproxy")
	must(err)
	fmt.Printf("exists: %+v %+v\n", exists, stat)

	if(! exists) {
		grproxy, err := conn.Create("/grproxy", []byte("grproxy:80"), flags, acl)
		must(err)
		fmt.Printf("create: %+v\n", grproxy)
	}

	childchn, errors := monitorGserver(conn,"/grproxy")

	go func() {
		for {
			select {

			case children := <-childchn:
				fmt.Printf("Show all children: %+v\n", children)
				var temp []string
				for _, child := range children {
					gserve_urls, _, err := conn.Get("/grproxy/" + child)
					temp = append(temp, string(gserve_urls))
					if(err != nil) {
						fmt.Printf("from child error: %+v\n", err)
					}			
				}
			urls = temp
			fmt.Printf(" All gserver urls: %+v\n", urls)
			case err := <-errors:
				fmt.Printf("channel error:  %+v\n", err)
			}
		}
	}()

	proxy := NewMultipleHostReverseProxy()
	log.Fatal(http.ListenAndServe(":9097", proxy))
	
}




