package main

import (
	"log"
	"strconv"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	rabbitMQ, err := NewRabbitMQService("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	emailService := NewEmailService("sandbox.smtp.mailtrap.io", "2525", "8483b0b9fab2f3", "7c95e08b62e4e9")

	schedulerHandler := NewSchedulerHandler(rabbitMQ)

	r := gin.Default()
	r.POST("/schedule", schedulerHandler.ScheduleEmail)

	go consumeMessages(rabbitMQ, emailService)

	r.Run(":8000")
}

func consumeMessages(rabbitMQ *RabbitMQService, emailService *EmailService) {
	queueName := "scheduled_emails"

	// Declare the queue
	_, err := rabbitMQ.DeclareQueue(queueName)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := rabbitMQ.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	log.Println("Started consuming messages...")

	for msg := range msgs {
		var email Email
		if err := json.Unmarshal(msg.Body, &email); err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}

		if time.Now().After(email.SendAt) {
			if err := emailService.SendEmail(email); err != nil {
				log.Printf("Error sending email: %v", err)
			}
		} else {
			// Requeue the message with updated expiration
			body, _ := json.Marshal(email)
			delay := email.SendAt.Sub(time.Now()).Milliseconds()
			if delay < 0 {
				delay = 0
			}
			expirationStr := strconv.FormatInt(delay, 10)

			err := rabbitMQ.Publish(
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
				log.Printf("Error requeueing message: %v", err)
			}
		}
	}
}
