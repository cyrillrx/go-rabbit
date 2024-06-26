package rabbitmq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel, exchangeName string) error {

	return ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

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

func sendMessage(ch *amqp.Channel, exchangeName, routingKey, msg string) error {

	return ch.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
}

// Checks if the AMQP connection is still on
func checkConnection(
	log Logger,
	conn *amqp.Connection,
	foreverChan chan bool,
	interval time.Duration) {

	for {
		tryToConnect(log, conn, foreverChan)
		time.Sleep(interval)
	}
}

// Checks if the AMQP connection is still on
func tryToConnect(
	log Logger,
	conn *amqp.Connection,
	foreverChan chan bool) {

	log.Info("Checking connection to Rabbit")
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("Checking channel failed: %s", err)
		close(foreverChan)
		return
	}
	ch.Close()
}
