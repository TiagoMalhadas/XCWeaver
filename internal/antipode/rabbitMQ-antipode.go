package antipode

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	connection *amqp.Connection
	queue      string
}

// How can I close the connection?
// É má prática manter a conexão sempre aberta?
func CreateRabbitMQ(rabbit_host string, rabbit_port string, rabbit_user string, rabbit_password string, queue string) RabbitMQ {

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")
	if err != nil {
		fmt.Println(err)
		return RabbitMQ{}
	}
	//defer conn.Close()

	return RabbitMQ{conn, queue}
}

func (r RabbitMQ) write(ctx context.Context, _ string, key string, obj AntiObj) error {

	jsonAntiObj, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	channel, err := r.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		r.queue, // Queue name
		false,   // Durable
		false,   // Delete when unused
		false,   // Exclusive
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		return err
	}

	err = channel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonAntiObj,
		})
	if err != nil {
		return err
	}

	return err
}

func (r RabbitMQ) read(ctx context.Context, _ string, key string) (AntiObj, error) {

	channel, err := r.connection.Channel()
	if err != nil {
		return AntiObj{}, err
	}
	//defer channel.Close()

	queue, err := channel.QueueDeclare(
		r.queue, // Queue name
		false,   // Durable
		false,   // Delete when unused
		false,   // Exclusive
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		return AntiObj{}, err
	}

	// Consume one message from the queue
	err = channel.Qos(1, 0, false)
	if err != nil {
		return AntiObj{}, err
	}
	msgs, err := channel.Consume(
		queue.Name, // Queue
		"",         // Consumer tag
		false,      // Auto-ack
		false,      // Exclusive
		false,      // No local
		false,      // No wait
		nil,        // Args
	)
	if err != nil {
		log.Fatalf("Failed to consume messages from queue: %v", err)
	}

	for len(msgs) == 0 {

	}
	channel.Close()

	// Wait for the first message to arrive and send an acknowledgement
	msg := <-msgs
	err = msg.Ack(true)
	if err != nil {
		return AntiObj{}, err
	}

	var antiObj AntiObj

	err = json.Unmarshal(msg.Body, &antiObj)
	if err != nil {
		return AntiObj{}, fmt.Errorf("Failed to unmarshal JSON: %v", err)
	}

	return antiObj, err
}

func (r RabbitMQ) barrier(ctx context.Context, lineage []WriteIdentifier, datastoreID string) error {
	return nil
}
