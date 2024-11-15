package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var log = logrus.New()

func main() {
	log.WithFields(logrus.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server CORE start")

	go HTTPAPIServer()
	go RTSPServer()
	go Storage.StreamChannelRunAll()

	// Thêm chức năng đọc từ RabbitMQ
	go RabbitMQConsumer()

	signalChanel := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChanel
		log.WithFields(logrus.Fields{
			"module": "main",
			"func":   "main",
		}).Info("Server receive signal", sig)
		done <- true
	}()

	log.WithFields(logrus.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server start success and waiting for signals")
	<-done
	Storage.StopAll()
	time.Sleep(2 * time.Second)
	log.WithFields(logrus.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server stop working by signal")
}

// RabbitMQConsumer handles consuming messages from RabbitMQ
func RabbitMQConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "rabbitmq",
			"func":   "RabbitMQConsumer",
		}).Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "rabbitmq",
			"func":   "RabbitMQConsumer",
		}).Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := "vms_queue"
	queue, err := ch.QueueDeclare(
		queueName, // Tên Queue
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "rabbitmq",
			"func":   "RabbitMQConsumer",
		}).Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		queue.Name, // Queue name
		"",         // Consumer name
		true,       // Auto-ack
		false,      // Exclusive
		false,      // No-local
		false,      // No-wait
		nil,        // Arguments
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "rabbitmq",
			"func":   "RabbitMQConsumer",
		}).Fatalf("Failed to register a consumer: %v", err)
	}

	log.WithFields(logrus.Fields{
		"module": "rabbitmq",
		"func":   "RabbitMQConsumer",
	}).Info("RabbitMQ consumer started and waiting for messages")

	for msg := range msgs {
		log.WithFields(logrus.Fields{
			"module": "rabbitmq",
			"func":   "RabbitMQConsumer",
		}).Infof("Received message: %s", msg.Body)

		// Xử lý message ở đây
		processMessage(msg.Body)
	}
}

func processMessage(body []byte) {
	log.WithFields(logrus.Fields{
		"module": "rabbitmq",
		"func":   "processMessage",
	}).Infof("Processing message: %s", body)

}
