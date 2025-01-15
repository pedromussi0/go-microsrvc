package main

import (
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	//start listening for messages
	log.Println("Listening for messages")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn, "events")
	if err != nil {
		panic(err)
	}
	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	// connect to rabbitmq
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//dont continue until rabbit is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("Failed to connect to RabbitMQ. Retrying in 1 second")
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if counts > 5 {
			log.Println("Failed to connect to RabbitMQ after 5 retries")
			break
		}

		backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
		log.Println("Retrying in ", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
