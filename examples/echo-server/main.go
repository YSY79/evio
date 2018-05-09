package main

import (
	"github.com/YSY79/evio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"bytes"
	"time"
	"strings"
	"os"
)

func main() {
	//var events evio.Events
	//events.Data = func(id int, in []byte) (out []byte, action evio.Action) {
	//
	//	data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(in), simplifiedchinese.GB18030.NewDecoder()))
	//	fmt.Println("长度：",len(data))
	//	fmt.Print("内容：",string(data))
	//	out = in
	//	return
	//}
	//events.Serving = func(srvin evio.Server) (action evio.Action) {
	//	return
	//}
	//events.Opened= func(id int, info evio.Info) (out []byte, opts evio.Options, action evio.Action) {
	// c,_:=evio.AllNetConns.Load(id)
	//   if(c!=nil){
	//	   c.(*evio.NetConn).Write([]byte("hello client!"))
	//
	//   }
	//   fmt.Println(id)
	//   fmt.Println(info)
	//	return
	//}
	//if err := evio.Serve(events, "tcp://:5000?reuseport=false"); err != nil {
	//	panic(err.Error())
	//}
	testServe("tcp-net", ":9992", false, 10)
}

func testServe(network, addr string, unix bool, nclients int) {
	//var started bool
	var connected int
	var disconnected int
	socket := strings.Replace(addr, ":", "socket", 1)
	fmt.Println(socket)
	var events evio.Events
	events.Serving = func(srv evio.Server) (action evio.Action) {
		return
	}
	events.Opened = func(id int, info evio.Info) (out []byte, opts evio.Options, action evio.Action) {
		connected++
		out = []byte("sweetness\r\n")
		opts.TCPKeepAlive = time.Minute * 5
		//c,_:=evio.AllNetConns.Load(id)
		//c.(*evio.NetConn).Write([]byte("hello!client"))
		if info.LocalAddr == nil {
			panic("nil local addr")
		}
		if info.RemoteAddr == nil {
			panic("nil local addr")
		}
		return
	}
	events.Closed = func(id int, err error) (action evio.Action) {
		disconnected++
		if connected == disconnected && disconnected == nclients {
			action = evio.Shutdown
		}
		return
	}
	events.Data = func(id int, in []byte) (out []byte, action evio.Action) {
			data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(in), simplifiedchinese.GB18030.NewDecoder()))
			fmt.Println("长度：",len(data))
			fmt.Print("内容：",string(data))
			out = in
			return
	}
	//events.Tick = func() (delay time.Duration, action evio.Action) {
	//	if !started {
	//		for i := 0; i < nclients; i++ {
	//			go startClient(network, addr)
	//		}
	//		started = true
	//	}
	//	delay = time.Second / 5
	//	return
	//}


	var err error
	if unix {
		socket := strings.Replace(addr, ":", "socket", 1)
		os.RemoveAll(socket)
		defer os.RemoveAll(socket)
		err = evio.Serve(events, network+"://"+addr, "unix://"+socket)
	} else {
		err = evio.Serve(events, network+"://"+addr)
	}


	fmt.Println(network+"://"+addr)

	if err != nil {
		panic(err)
	}
}
