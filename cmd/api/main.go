package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/alireza-aliabadi/golang-stocking-update/internal/config"
	"github.com/alireza-aliabadi/golang-stocking-update/internal/rabbitmq"
	stockHttp "github.com/alireza-aliabadi/golang-stocking-update/internal/stock/delivery/http"
)

func main() {
	rabbitUrl := config.LoadConf().RabbitmqUrl
	rabbitmqConn, err := rabbitmq.Connect(rabbitUrl)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer rabbitmqConn.Close()

	e := echo.New()

	handler := stockHttp.NewStockHandler(rabbitmqConn)
	handler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}