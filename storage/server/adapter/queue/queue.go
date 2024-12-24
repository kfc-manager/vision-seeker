package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Close() error
	Push(msg []byte) error
}

type queue struct {
	conn    *amqp.Connection
	name    string
	channel *amqp.Channel
}

func New(host, port, name string) (*queue, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port))
	if err != nil {
		return nil, err
	}

	prod, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = prod.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &queue{
		conn:    conn,
		channel: channel,
		name:    name,
	}, nil
}

func (q *queue) Close() error {
	chanErr := q.channel.Close()
	connErr := q.conn.Close()

	if chanErr != nil {
		return chanErr
	}
	if connErr != nil {
		return connErr
	}

	return nil
}

func (q *queue) Push(msg []byte) error {
	return q.channel.PublishWithContext(
		context.Background(),
		"",     // exchange
		q.name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "test/plain; charset=UTF-8",
			Body:        msg,
		})
}
