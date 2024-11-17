package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vishal/microservice-kafka-golang/types"
)

const websocketURL string = "ws://localhost:3000/ws"

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		coordinatesArr := []types.OnboardUnit{}
		for i := 0; i < 20; i++ {
			lat, long := GenCoordinates()
			obuUnit := types.OnboardUnit{
				OID:  uuid.New().String(),
				Lat:  lat,
				Long: long,
			}
			coordinatesArr = append(coordinatesArr, obuUnit)
			if err = conn.WriteJSON(coordinatesArr); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(time.Second * 5)
		fmt.Printf("Coordinates are %+v \n", coordinatesArr)
	}
}

func GenerateRandomNumber() float32 {
	randomInteger := rand.Intn(100 + 1)
	randomFloat := rand.Float32()
	return float32(randomInteger) + randomFloat
}

func GenCoordinates() (float32, float32) {

	lat := GenerateRandomNumber()
	long := GenerateRandomNumber()
	return lat, long
}

//	func generateOBUIds(n int) {
//		ids := make([]int, n)
//		for i := 0; i < n; i++ {
//			ids[i] = rand.Intn(math.MaxInt)
//		}
//		fmt.Println(ids)
//	}
