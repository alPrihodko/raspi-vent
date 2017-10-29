package main

import (
	"bytes"
	"encoding/json"
	"log"
	"raspi-vent/wsHandler"

	"sync"
	"time"

	"golang.org/x/net/websocket"
)

/*
MsgCommand - simple comm protocol to deliver infor to server
*/
type MsgCommand struct {
	sync.Mutex
	Command string `json:"Command, string"`
	Object  string `json:"Object, string"`
	Param1  string `json:"Param1, string"`
	Param2  string `json:"Param2, string"`
}

func execute(c *MsgCommand) error {
	defer c.Unlock()
	c.Lock()

	return nil
}

func relHandler(ws *websocket.Conn) {

	var id = int32(time.Now().Unix())

	wh = wsHandler.New(id, ws)

	defer func() {
		wh.Destroy(id)
	}()

	msg := make([]byte, 512)
	for {
		n, err := ws.Read(msg)
		if err != nil {
			log.Println(err)
			return
		}

		x := new(MsgCommand)
		log.Printf("Receive: %s\n", msg[:n])
		if err := json.NewDecoder(bytes.NewReader(msg)).Decode(x); err == nil {
			execute(x)
		} else {
			log.Println(err)
		}
	}
}
