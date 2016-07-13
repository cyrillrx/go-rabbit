package rabbitmq

import (
	"time"

	"github.com/rs/xlog"
	"github.com/streadway/amqp"
)

func declareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {

	return ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func bindQueue(ch *amqp.Channel, q amqp.Queue, exchangeName, routingKey string) error {

	return ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
}

func consumeQueue(ch *amqp.Channel, q amqp.Queue) (<-chan amqp.Delivery, error) {

	return ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

// Checks if the AMQP connection is still on
func checkConnection(
	log xlog.Logger,
	conn *amqp.Connection,
	foreverChan chan bool,
	interval time.Duration) {

	for {
		ch, err := conn.Channel()
		if err != nil {
			log.Errorf("Checking channel failed: %s", err)
			close(foreverChan)
			return
		}
		defer ch.Close()

		time.Sleep(interval)
	}
}
