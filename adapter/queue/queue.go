package queue

import (
	"context"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Close() error
	Push(msg string) error
	Pull() (string, error)
}

type queue struct {
	conn    *amqp.Connection
	name    string
	produce *amqp.Channel
	consume *amqp.Channel
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
	cons, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = prod.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-max-length": 1000000, // Maximum messages
			"x-overflow":   "reject-publish",
		},
	)
	if err != nil {
		return nil, err
	}

	return &queue{
		conn:    conn,
		produce: prod,
		consume: cons,
		name:    name,
	}, nil
}

func (q *queue) Close() error {
	prodErr := q.produce.Close()
	consErr := q.produce.Close()
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

func (q *queue) Push(msg string) error {
	err := q.produce.PublishWithContext(
		context.Background(),
		"",     // exchange
		q.name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		return err
	}

	return nil
}

func (q *queue) Pull() (string, error) {
	del, ok, err := q.consume.Get(
		q.name, // queue name
		true,   // auto acknowledge
	)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("no message left in queue")
	}

	return string(del.Body), nil
}
