package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var status struct {
	Ok         bool      `json:"ok"`
	Msg        string    `json:"msg"`
	LastUpdate time.Time `json:"last_update"`
}

var serveAddr = flag.String("addr", ":80", "the address to listen on")
var updateFreqMinutes = flag.Duration("duration", time.Minute*1, "the duration between two detects, in minutes")

func init() {
	flag.Parse()
}

func main() {
	checkConfiguration()
	initStatus()
	go startCheckerServer()
	startHttpServer()
}

func initStatus() {
	status.Ok = true
	status.Msg = "监控探针已经启动"
	status.LastUpdate = time.Now()
}

func checkConfiguration() {
	if *updateFreqMinutes >= 365*24*60*time.Minute {
		log.Fatalln("The specified duration is too long!")
	}
}

func startCheckerServer() {
	for now := range time.Tick(*updateFreqMinutes) {
		log.Printf("Checking network status.")
		checkHttpConnection()
		status.LastUpdate = now
	}
}

func checkHttpConnection() {
	resp, err := http.Get("http://www.baidu.com/")
	if err == nil {
		if resp.StatusCode == 200 {
			status.Ok = true
			status.Msg = "一切正常"
		} else {
			status.Ok = false
			status.Msg = "似乎不能正常连接到百度"
		}
	} else {
		if strings.Contains(err.Error(), ": EOF") {
			status.Ok = false
			status.Msg = "外网连接被网络中心阻断"
		} else {
			status.Ok = false
			status.Msg = fmt.Sprintf("访问失败：%s", err.Error())
		}
	}
	log.Println(status)
}

func startHttpServer() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/api/status", func(writer http.ResponseWriter, request *http.Request) {
		encoder := json.NewEncoder(writer)
		_ = encoder.Encode(status)
	})

	if http.ListenAndServe(*serveAddr, nil) != nil {
		log.Fatalf("Unable to listen on address `%s`.\n", *serveAddr)
	}
}
