package worker

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	stock "github.com/alireza-aliabadi/golang-stocking-update/internal/models"
	rabbitmq "github.com/alireza-aliabadi/golang-stocking-update/internal/rabbitmq"
	stockService "github.com/alireza-aliabadi/golang-stocking-update/internal/stock/service"
)

type Consumer struct {
	Client *rabbitmq.RabbitClient
	Service *stockService.StoreUpdater
}

func NewConsumer(client *rabbitmq.RabbitClient, svc *stockService.StoreUpdater) *Consumer {
	return &Consumer{
		Client: client,
		Service: svc,
	}
}

func (c *Consumer) processMessages(id int, messages <-chan amqp.Delivery) {
	for d := range messages {
		var payload stock.StockUpdatePayload
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			log.Printf("[Worker %d] Invalid JSON: %v", id, err)
			d.Nack(false, false)
			continue
		}

		err := c.Service.UpdateAllStores(payload)

		if err != nil {
			log.Printf("[Worker %d] Failed: %v", id, err)
			d.Nack(false, true)
		} else {
			d.Ack(false)
		}
	}
}

func (c *Consumer) Start(workerCount int) {
	err := c.Client.Channel.Qos(workerCount*2, 0, false)
	if err != nil {
		log.Fatal("Failed to set QoS: %v", err)
	}

	messages, err := c.Client.Channel.Consume(
		c.Client.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan bool)

	for i := 0; i < workerCount; i++ {
		go c.processMessages(i, messages)
	}

	log.Printf("Worker pool started with %d workers waiting for messages...", workerCount)
	<-forever
}