package main

import (
	"github.com/YSY79/evio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"bytes"
)

func main() {
	var events evio.Events
	events.Data = func(id int, in []byte) (out []byte, action evio.Action) {

		data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(in), simplifiedchinese.GB18030.NewDecoder()))
		fmt.Println("长度：",len(data))
		fmt.Print("内容：",string(data))
		out = in
		return
	}
	events.Opened= func(id int, info evio.Info) (out []byte, opts evio.Options, action evio.Action) {
	 c,_:=evio.AllConnections.Load(id)
		c.(*evio.NetConn).Write([]byte("hello client!"))
		return
	}
	if err := evio.Serve(events, "tcp://localhost:5000"); err != nil {
		panic(err.Error())
	}
}
