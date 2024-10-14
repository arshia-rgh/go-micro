package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

var exchangeName = os.Getenv("RABBITMQ_EXCHANGE")

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

	return declareExchange(exchangeName, channel)

}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	q, err := declareRandomQueue(channel)

	if err != nil {
		return err
	}

	for _, v := range topics {
		err = channel.QueueBind(
			q.Name,
			v,
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

}
