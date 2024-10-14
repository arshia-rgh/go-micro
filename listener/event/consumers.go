package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()

	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	exchangeName := os.Getenv("RABBITMQ_EXCHANGE")
	return declareExchange(exchangeName, channel)

}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
