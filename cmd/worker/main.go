// runs background processor
package main

import (
	"log"

	"github.com/alireza-aliabadi/golang-stocking-update/internal/config"
	"github.com/alireza-aliabadi/golang-stocking-update/internal/rabbitmq"
	"github.com/alireza-aliabadi/golang-stocking-update/internal/stock/service"
	"github.com/alireza-aliabadi/golang-stocking-update/internal/stock/worker"
)

func main() {
	rabbitmqUrl := config.LoadConf().RabbitmqUrl
	rabbitmqConn, err := rabbitmq.Connect(rabbitmqUrl)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer rabbitmqConn.Close()

	svc := service.NewStoreUpdater()

	consumer := worker.NewConsumer(rabbitmqConn, svc)

	consumer.Start(5)
}