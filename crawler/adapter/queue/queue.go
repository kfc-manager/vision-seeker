package queue

import (
	"context"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Close() error
	Push(msg []byte) error
	Pull() ([]byte, error)
}

type queue struct {
	conn *amqp.Connection
	name string
	prod *amqp.Channel
	cons *amqp.Channel
	msgs <-chan amqp.Delivery
}

func New(host, port, name string, maxSize int) (*queue, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port))
	if err != nil {
		return nil, err
	}

	prod, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	var args amqp.Table = nil
	if maxSize > 0 {
		args = amqp.Table{
			"x-max-length": maxSize, // Maximum messages
			"x-overflow":   "reject-publish",
		}
	}
	_, err = prod.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,
	)
	if err != nil {
		return nil, err
	}

	cons, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	msgs, err := cons.Consume(
		name,  // queue name
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &queue{
		conn: conn,
		prod: prod,
		cons: cons,
		msgs: msgs,
		name: name,
	}, nil
}

func (q *queue) Close() error {
	prodErr := q.prod.Close()
	consErr := q.cons.Close()
	connErr := q.conn.Close()

	if prodErr != nil {
		return prodErr
	}
	if consErr != nil {
		return consErr
	}
	if connErr != nil {
		return connErr
	}

	return nil
}

func (q *queue) Push(msg []byte) error {
	return q.prod.PublishWithContext(
		context.Background(),
		"",     // exchange
		q.name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
}

func (q *queue) Pull() ([]byte, error) {
	msg, ok := <-q.msgs
	if !ok {
		return nil, errors.New("queue has been closed")
	}
	return msg.Body, nil
}
