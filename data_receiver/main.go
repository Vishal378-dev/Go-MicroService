package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/vishal/microservice-kafka-golang/types"
)

//https://www.tencentcloud.com/document/product/597/60360

var kafkaCoordinateTopic = "CoordinateData"

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWs)
	fmt.Println("Data receiver")
	http.ListenAndServe(":3000", nil)

}

type DatReceiver struct {
	msg      chan []types.OnboardUnit
	conn     *websocket.Conn
	producer *kafka.Producer
}

func NewDataReceiver() (*DatReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &DatReceiver{
		msg:      make(chan []types.OnboardUnit, 128),
		producer: p,
	}, nil
}

func (dr *DatReceiver) produceData(data []types.OnboardUnit) error {
	MarshaledData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = dr.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaCoordinateTopic,
			Partition: kafka.PartitionAny},
		Value: MarshaledData,
	}, nil)
	return err
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
		// dr.msg <- data

		if err := dr.produceData(data); err != nil {
			fmt.Println("Kafka err- ", err)
		}
	}
}
