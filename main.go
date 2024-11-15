package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Hello there i am using ubuntu now")
	for {
		lat, long := GenCoordinates()
		fmt.Printf("latitude is %f and longitube is %f \n", lat, long)
		time.Sleep(time.Second * 1)
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
