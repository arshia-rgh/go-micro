package event

import amqp "github.com/rabbitmq/amqp091-go"

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
