package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QueueName = "notifications"

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}
type NotificationMessage struct {
	ID string `json:"id"`
	Channel string `json:"channel"`
	Recipient string `json:"recipient"`
	Subject string `json:"subject"`
	Body string `json:"body"`
	SendAt string `json:"send_at"`
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial:%w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq channel:%w", err)
	}

	_, err = ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq queue declare:%w", err)
	}

	log.Println("RabbitMQ connected")
	return &RabbitMQ{
		conn: conn,
		ch:   ch,
		}, nil
}

func (r *RabbitMQ) Close() {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitMQ) Publish(ctx context.Context, msg NotificationMessage) error {
	body,err:= json.Marshal(msg)
	if err!=nil{
		return fmt.Errorf("marshal message: %w", err)
	}

	err= r.ch.PublishWithContext(ctx,
		"",
		QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",	
			Body: body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("publish message: %w", err)
	}
	log.Printf("rabbitmq send message by id=%s", msg.ID)
	return nil
}
