package main

import (
	"fmt"
	"os"
	"log"
	"github.com/samuel/go-zookeeper/zk"
    "net/http"
	"strings"
	"time"
	"bytes"
	"encoding/json"
	"io/ioutil"
)

var zookeeper string = "zookeeper" 
var hbase_host string = "hbase" 
var server_name string = "Unknown"

func must(err error) {
	if err != nil {
		//panic(err)
		fmt.Printf("%+v From must \n", err)
	}
}

func connect() *zk.Conn {
	zksStr := zookeeper+":2181" 
	zks := strings.Split(zksStr, ",")
	conn, _, err := zk.Connect(zks, time.Second)
	must(err)
	return conn
}

func encoder(unencodedJSON []byte) string {
	// get go object from json byte 
	var unencodedRows RowsType
	json.Unmarshal(unencodedJSON, &unencodedRows)

	//  encode all fields value of go object , return EncRowsType
	encodedRows := unencodedRows.encode()

	// convert to json byte[] from go object (EncRowsType)
	encodedJSON, _ := json.Marshal(encodedRows)

	return string(encodedJSON)
}

func decoder(encodedJSON []byte) string {

	// get go object from json byte 
	var encodedRows EncRowsType
	fmt.Println("From decoder test print: ", string(encodedJSON))
	json.Unmarshal(encodedJSON, &encodedRows)
	fmt.Println("From decoder first: ", encodedRows)

	//  decode all fields value of go object , return RowsType
	decodedRows, err := encodedRows.decode()
	if err != nil {
		fmt.Println("%+v", err)
	}
    fmt.Println("From decoder second: ", decodedRows		)
	// convert to json byte[] from go object (RowsType)
	deCodedJSON,_:= json.Marshal(decodedRows)

	//fmt.Println("From decoder method: ", string(deCodedJSON))
	return string(deCodedJSON)
}

//func postToHbase(writer http.ResponseWriter, req *http.Request) {
func postToHbase(encodedJSON string) {

	/* 
	get []byte format data from request body  
	 encodedJsonByte,err := ioutil.ReadAll(req.Body)
	 must(err)

	 get encoded data from []byte type
	 encodedJSON := encoder(encodedJsonByte)	
     fmt.Println("encodedJSON :  %+v\n", string(encodedJSON));
	 */

	req_url := hbase_host + "/se2:library/fakerow"
	
	resp, err := http.Post( req_url, "application/json" ,bytes.NewBuffer([]byte(encodedJSON)))

	if err != nil {
		fmt.Println("Error from response: %+v", err)
		return	
	}

	fmt.Println("Post Response: ", resp.Status)
	defer resp.Body.Close()
}

//func getFromHbase(writer http.ResponseWriter, req *http.Request) {
func getFromHbase() string {

	req_url := hbase_host + "/se2:library/*"

	// resp, getErr := http.Get(req_url)	
	req, _ := http.NewRequest("GET", req_url, nil)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, getErr := client.Do(req)
	must(getErr)

	fmt.Println("Get Response: ", resp.Status)
	
	encodedJsonByte,err := ioutil.ReadAll(resp.Body)
	must(err)

	decodedJSON := decoder(encodedJsonByte)	
	//fmt.Println("decodedJSON :  ", string(decodedJSON))
	defer resp.Body.Close()
	return decodedJSON
}
 
func handler(writer http.ResponseWriter, req *http.Request) {

	if (req.Method == "POST" || req.Method == "PUT") {

		encodedJsonByte,err := ioutil.ReadAll(req.Body)
		must(err)

		// get encoded data from []byte type
		encodedJSON := encoder(encodedJsonByte)	
    	fmt.Println("encodedJSON : ", string(encodedJSON));

		req.Header.Set("Content-type", "application/json")
		postToHbase(encodedJSON)
		fmt.Fprintf(writer, "an %s\n", "POST")

	}else if(req.Method == "GET") {
		fmt.Printf("hello get")
		req.Header.Set("Accept", "application/json")
		responseData := getFromHbase()
		// TODO: create html using this data 
		// write it to ResponseWriter
		fmt.Fprintf(writer, "Response from hbase:\n\n %s\n", string(responseData))

	}else {
		fmt.Fprintf(writer, "Invalid Request from Client")
	}

	fmt.Fprintf(writer, "Proudly server by %s", server_name)
	
}

func startServer() {
    	http.HandleFunc("/library", handler)
    	log.Fatal(http.ListenAndServe(":9091", nil))
}

func main() {

	server_name = os.Getenv("servername")
	println("Starting server %s.......", server_name)

	conn := connect()
	defer conn.Close()

	flags := int32(zk.FlagEphemeral)
	acl := zk.WorldACL(zk.PermAll)

	// server_name:9091
	gserv, err := conn.Create("/grproxy/" + server_name, []byte( server_name+":9091"), flags, acl)
	must(err)
	fmt.Printf("create ephemeral node: %+v\n", gserv)
	
	
	startServer()
	//getFromHbase()

	//temp := []byte("{\"Row\":[{\"key\":\"test\",\"Cell\":[{\"column\":\"document:Chapter 2\",\"$\":\"value:this is test 2\"},{\"column\":\"metadata:Author\",\"$\":\"value:test auther!\"}]}]}")
	//postToHbase(encoder(temp))
	//temp1 := getFromHbase()	
	//println(temp1)

}



