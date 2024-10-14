package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (emitter *Emitter) setup() error {
	channel, err := emitter.connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	return declareExchange(exchangeName, channel)
}

func (emitter *Emitter) Push(event, key string) error {
	channel, err := emitter.connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	err = channel.Publish(
		exchangeName,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	return err
}

func NewEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{connection: conn}

	err := emitter.setup()

	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
