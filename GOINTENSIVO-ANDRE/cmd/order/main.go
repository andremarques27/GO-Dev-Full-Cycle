package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/devfullcycle/go-intensivo-andre/internal/infra/database"
	"github.com/devfullcycle/go-intensivo-andre/internal/usercase"
	"github.com/devfullcycle/go-intensivo-andre/pkg/rabbitmq"
	_ "github.com/mattn/go-sqlite3"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	orderRepository := database.NewOrderRepository(db)
	uc := usercase.NewCalculateFinalPrice(orderRepository)
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	msgRabbitmqChanel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChanel)
	rabbitmqWorker(msgRabbitmqChanel, uc)
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc *usercase.CalculateFinalPrice) {
	fmt.Println("Starting rabbitmq")
	for msg := range msgChan {
		var input usercase.OrderInput
		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			panic(err)
		}
		output, err := uc.Execute(input)
		if err != nil {
			panic(err)
		}
		msg.Ack(false)
		fmt.Println("Mensagem processada e salva no banco: ", output)
	}
}