package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
	"github.com/vishal/microservice-kafka-golang/types"
)

//https://www.tencentcloud.com/document/product/597/60360

var kafkaCoordinateTopic = "CoordinateData"

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

func ConnectProducer(brokerUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = time.Second

	prod, err := sarama.NewSyncProducer(brokerUrl, config)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (dr *DatReceiver) produceData(topic string, data []types.OnboardUnit) error {
	MarshaledData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	brokers := []string{"localhost:29092"}
	prod, err := ConnectProducer(brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer prod.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(MarshaledData),
	}
	// send message
	partition, offset, err := prod.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Message is stored in topic(%s) /partition (%d) / offset(%d)\n", topic, partition, offset)
	return nil
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
		// fmt.Printf("coordinates array is - %+v\n", data)
		err := dr.produceData(kafkaCoordinateTopic, data)
		if err != nil {
			log.Fatal(err)
		}
		// dr.msg <- data
	}
}

/*
parse request body into order
convert body into bytes
send the bytes to kafka
respond back to user
*/
