package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	Connection *amqp.Connection
	Channel *amqp.Channel
	Queue amqp.Queue
}

func Connect(url string) (*RabbitClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		"stock_updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		Connection: conn,
		Channel: ch,
		Queue: queue,
	}, nil
}

func (rc *RabbitClient) Close() {
	rc.Channel.Close()
	rc.Connection.Close()
}

func (rc *RabbitClient) Publish(payload interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	return rc.Channel.PublishWithContext(
		ctx,
		"",
		rc.Queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "application/json",
			Body: body,
		},
	)
}