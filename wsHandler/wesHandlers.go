package wsHandler

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type stringReport struct {
	Type  string
	Key   string
	Value string
}

/*
SocketConns - socket connections
*/
type SocketConns struct {
	ws   map[int32]*websocket.Conn
	lock *sync.Mutex
}

var conns SocketConns

type wsHandler struct {
	Wsid int32
}

/*WsHandler contains user preferences*/
type WsHandler wsHandler

func init() {
	conns = SocketConns{make(map[int32]*websocket.Conn), &sync.Mutex{}}
}

/*New - registers new connection*/
func New(id int32, ws *websocket.Conn) WsHandler {
	defer conns.lock.Unlock()
	conns.lock.Lock()
	conns.ws[id] = ws

	return WsHandler{id}
}

/*Destroy - destroys connection*/
func Destroy(id int32) error {
	defer conns.lock.Unlock()
	log.Println("Destroying ws connection")
	conns.lock.Lock()
	if _, ok := conns.ws[id]; ok {
		delete(conns.ws, id)
		log.Println("Removing connection id: ", id)
		return nil
	}

	return errors.New("Unable to Destroy socket handler")
}

/*Destroy - destroys connection*/
func (w WsHandler) Destroy(id int32) error {
	defer conns.lock.Unlock()
	log.Println("Destroying ws connection")
	conns.lock.Lock()
	if _, ok := conns.ws[id]; ok {
		delete(conns.ws, id)
		log.Println("Removing connection id: ", id)
		return nil
	}

	return errors.New("Unable to Destroy socket handler")
}

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

/*ReportWsEvent - reports to ws*/
func (w WsHandler) ReportWsEvent(evt string, st string) error {

	r := stringReport{evt, "state", st}

	b, err01 := json.Marshal(r)
	if err01 != nil {
		return err01
	}

	for _, ws := range conns.ws {
		m, err02 := ws.Write(b)
		if err02 != nil {
			return err02
		}
		log.Println(m)
	}

	return nil
}

/*
ReportData - reporting stream
*/
func ReportData(d []byte) error {

	for _, ws := range conns.ws {
		_, err02 := ws.Write(d)
		if err02 != nil {
			return err02
		}
	}

	return nil
}
