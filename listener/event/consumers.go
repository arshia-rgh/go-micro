package event

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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

	messages, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	log.Printf("waiting for message [Exchange, Queue] = [%v, %v]\n", exchangeName, q.Name)
	<-forever

	return nil

}

func handlePayload(payload Payload) error {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)

		if err != nil {
			log.Println(err)
		}

	case "auth":
		// auth

	default:
		err := logEvent(payload)

		if err != nil {
			log.Println(err)
		}

	}

}

func logEvent(entry Payload) error {

}
