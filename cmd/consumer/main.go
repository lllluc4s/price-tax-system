package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lllluc4s/go-fullcycle/internal/infra/database"
	"github.com/lllluc4s/go-fullcycle/internal/usecase"
	"github.com/lllluc4s/go-fullcycle/pkg/kafka"
	"github.com/lllluc4s/go-fullcycle/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./orders.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)
	usecase := usecase.CalculateFinalPrice{OrderRepository: repository}
	msgChanKafka := make(chan *ckafka.Message)
	topics := []string{"orders"}
	servers := "host.docker.internal:9094"

	fmt.Println("Kafka consumer has started!")
	go kafka.Consume(topics, servers, msgChanKafka)
	go kafkaWorker(msgChanKafka, usecase)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgRabbitmqChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel)
	rabbitmqWorker(msgRabbitmqChannel, usecase)
}

func kafkaWorker(msgChan chan *ckafka.Message, uc usecase.CalculateFinalPrice) {
	fmt.Println("Kafka worker has started!")
	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Value, &OrderInputDTO)
		if err != nil {
			panic(err)
		}

		outputDTO, errr := uc.Execute(OrderInputDTO)
		if errr != nil {
			panic(errr)
		}

		fmt.Printf("Kafka has processed order %s\n", outputDTO.ID)
	}
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc usecase.CalculateFinalPrice) {
	fmt.Println("Rabbitmq worker has started!")
	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &OrderInputDTO)
		if err != nil {
			panic(err)
		}

		outputDTO, errr := uc.Execute(OrderInputDTO)
		if errr != nil {
			panic(errr)
		}

		msg.Ack(false)
		fmt.Printf("Rabbitmq has processed order %s\n", outputDTO.ID)
	}
}
