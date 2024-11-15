package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vishal/microservice-kafka-golang/types"
)

func main() {
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.handleWs)
	fmt.Println("Data receiver")
	http.ListenAndServe(":3000", nil)
}

type DatReceiver struct {
	msg  chan []types.OnboardUnit
	conn *websocket.Conn
}

func NewDataReceiver() *DatReceiver {
	return &DatReceiver{
		msg: make(chan []types.OnboardUnit, 128),
	}
}

func (dr *DatReceiver) handleWs(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		WriteBufferSize: 1028,
		ReadBufferSize:  1028,
	}

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReadLoop()
}

func (dr *DatReceiver) wsReadLoop() {
	fmt.Println("client connected")
	for {
		var data []types.OnboardUnit
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("Read Error :", err)
			continue
		}
		fmt.Printf("coordinates array is - %+v\n", data)
		dr.msg <- data
	}
}
