package event

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(name string, ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		name,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
}
