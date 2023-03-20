package api

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"time"
)

type MqService struct {
	conn *amqp.Connection
}

var (
	RabbitMqUrl = os.Getenv("RABBIT_MQ_URL")
	QueueName   = os.Getenv("RABBIT_MQ_NAME")
)

func TryDial(url string, tries int) (*amqp.Connection, error) {
	var lastErr error
	for i := 0; i < tries; i++ {
		conn, err := amqp.Dial(url)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		time.Sleep(5 * time.Second)
	}
	return nil, lastErr
}

func NewMqService() (*MqService, error) {
	conn, err := TryDial(RabbitMqUrl, 10)
	if err != nil {
		return nil, err
	}

	return &MqService{
		conn: conn,
	}, nil
}

func (service *MqService) Upload(req *RequestData) error {
	ch, err := service.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(QueueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	ctx := context.Background()
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	return err
}
