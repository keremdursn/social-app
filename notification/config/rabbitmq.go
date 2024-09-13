package config

import (
	"log"

	"github.com/streadway/amqp"
)

func ConnectRabbitmq() *amqp.Connection {
	connection, err := amqp.Dial("127.0.0.1:5672")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	// defer connection.Close()
	return connection
}

// Kanal oluştur ve kuyruk tanımla
func CreateQueue(ch *amqp.Channel, queueName string) amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // Queue ismi
		false,     // Durable (kalıcı)
		false,     // Delete when unused (boşsa sil)
		false,     // Exclusive (özel)
		false,     // No-wait (hemen oluştur)
		nil,       // Extra arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	return q
}

// Mesaj gönderme işlemi
func PublishMessage(ch *amqp.Channel, queueName, message string) {
	q := CreateQueue(ch, queueName)

	err := ch.Publish(
		"",     // Exchange
		q.Name, // Routing key (queue name)
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
	log.Printf("Message sent to queue %s: %s", q.Name, message)
}

func ConsumeMessages(ch *amqp.Channel, queueName string) {
	q := CreateQueue(ch, queueName)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("Waiting for messages in queue: %s", queueName)
}
