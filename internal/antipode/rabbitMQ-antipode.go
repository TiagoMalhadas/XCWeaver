package antipode

import (
	"context"
	"encoding/json"
	"fmt"

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
		return AntiObj{}, err
	}

	// Consume one message from the queue
	var ok bool
	var msg amqp.Delivery
	for !ok {
		msg, ok, err = channel.Get(queue.Name, false)
		if err != nil {
			return AntiObj{}, err
		}
	}

	var antiObj AntiObj
	err = json.Unmarshal(msg.Body, &antiObj)
	if err != nil {
		err = msg.Ack(true)
		if err != nil {
			return AntiObj{}, err
		}
		return AntiObj{}, fmt.Errorf("Failed to unmarshal JSON: %v", err)
	}

	err = msg.Ack(true)
	if err != nil {
		return AntiObj{}, err
	}

	return antiObj, err
}

func (r RabbitMQ) consume(ctx context.Context, _ string, stop chan struct{}) (<-chan AntiObj, error) {
	channel, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		r.queue, // Queue name
		false,   // Durable
		false,   // Delete when unused
		false,   // Exclusive
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}

	antipodeObjctsChan := make(chan AntiObj)

	go func() {
		defer close(antipodeObjctsChan)
		defer channel.Close()
		//requeue non-processed messages
		defer func(<-chan AntiObj) {
			for d := range antipodeObjctsChan {
				jsonAntiObj, err := json.Marshal(d)
				if err != nil {
					fmt.Errorf(err.Error())
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
					fmt.Errorf(err.Error())
				}
			}
		}(antipodeObjctsChan)

		select {
		//channel is closed
		case <-stop:
			return
		default:
			for d := range msgs {
				var antiObj AntiObj
				err := json.Unmarshal(d.Body, &antiObj)
				if err != nil {
					fmt.Errorf(err.Error())
				}
				antipodeObjctsChan <- antiObj
				d.Ack(true)
			}
		}
	}()

	return antipodeObjctsChan, nil
}

func (r RabbitMQ) barrier(ctx context.Context, lineage []WriteIdentifier, datastoreID string) error {
	return nil
}
