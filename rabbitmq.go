package main

import (
	"encoding/json"
	"time"

	"fmt"
	"strconv"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQService(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{conn: conn, channel: ch}, nil
}

func (s *RabbitMQService) ScheduleEmail(email Email) error {
	queueName := "scheduled_emails"
	_, err := s.DeclareQueue(queueName)
	if err != nil {
		return fmt.Errorf("Failed to declare queue: %v", err)
	}

	body, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("Failed to marshal email: %v", err)
	}

	delay := email.SendAt.Unix() - time.Now().Unix()
	if delay < 0 {
		delay = 0
	}

	expirationStr := strconv.FormatInt(delay, 10)

	err = s.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Expiration:  expirationStr,
		},
	)

	if err != nil {
		return fmt.Errorf("Failed to publish message: %v", err)
	}

	return nil
}

// Add this to services/rabbitmq.go

func (s *RabbitMQService) DeclareQueue(name string) (amqp.Queue, error) {
	return s.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (s *RabbitMQService) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return s.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (s *RabbitMQService) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return s.channel.Publish(exchange, key, mandatory, immediate, msg)
}
