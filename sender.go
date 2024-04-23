package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Sender connects to an AMQP server and performs queue and exchange bindings
type Sender struct {
	log          Logger
	conn         *amqp.Connection
	ch           *amqp.Channel
	deliveryChan <-chan amqp.Delivery
	handlers     []func(amqp.Delivery)

	amqpURL      string
	amqpLogURL   string // Obfuscated version of the amqpURL without credentials, suitable for logging
	exchangeName string
	routingKey   string

	tryConnect  int
	tryChannel  int
	tryExchange int
}

// NewSender Created a new AMQP sender
func NewSender(logger Logger, amqpURL, exchangeName, routingKey string) *Sender {

	s := &Sender{
		log:          logger,
		amqpURL:      amqpURL,
		amqpLogURL:   string(urlCredRe.ReplaceAll([]byte(amqpURL), []byte(":***@"))),
		exchangeName: exchangeName,
		routingKey:   routingKey,
		tryConnect:   0,
		tryChannel:   0,
	}

	s.setup()

	return s
}

func (s *Sender) setup() (err error) {

	// Connect to server
	s.tryConnect++
	//c.conn, err = amqp.Dial(c.amqpURL)
	s.conn, err = amqp.DialConfig(s.amqpURL, amqpConfig)
	if err != nil {
		return fmt.Errorf("Failed to connect to RabbitMQ: try %d, %v", s.tryConnect, err)
	}
	s.tryConnect = 0
	s.log.Infof("Sender connected to AMQP server at: %s", s.amqpLogURL)

	// Open channel
	s.tryChannel++
	s.ch, err = s.conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: try %d, %v", s.tryChannel, err)
	}
	s.tryChannel = 0

	// Declare exchange
	s.tryExchange++
	err = declareExchange(s.ch, s.exchangeName)
	if err != nil {
		return fmt.Errorf("Failed to declare a exchange: try %d, %s", s.tryExchange, err)
	}
	s.tryExchange = 0

	return err
}

// Publish Publishes an AMQP message
func (s *Sender) Publish(msg string) error {

	return sendMessage(s.ch, s.exchangeName, s.routingKey, msg)
}
