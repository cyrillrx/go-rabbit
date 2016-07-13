package rabbitmq

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/rs/xlog"
	"github.com/streadway/amqp"
)

const (
	retryDelay        = 1 * time.Minute
	checkConnInterval = 5 * time.Minute
)

var (
	urlCredRe = regexp.MustCompile(":.*?@")

	amqpConfig = amqp.Config{
		Heartbeat: 5 * time.Second,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 2*time.Second)
		},
	}
)

// Consumer connects to an AMQP server and performs queue and exchange bindings
type Consumer struct {
	log          xlog.Logger
	conn         *amqp.Connection
	ch           *amqp.Channel
	deliveryChan <-chan amqp.Delivery
	handlers     []func(amqp.Delivery)
	exchanges    map[string][]string // map of exchange names and corresponding routing keys

	amqpURL    string
	amqpLogURL string // Obfuscated version of the amqpURL without credentials (for logging)
	queueName  string

	tryConnect int
	tryChannel int
}

// NewConsumer Creates a RabbitMQ consumer for the given AMQP URL
func NewConsumer(logger xlog.Logger, amqpURL, queueName string, exchangeMap map[string][]string) *Consumer {

	return &Consumer{
		log:        logger,
		amqpURL:    amqpURL,
		amqpLogURL: string(urlCredRe.ReplaceAll([]byte(amqpURL), []byte(":***@"))),
		exchanges:  exchangeMap,
		queueName:  queueName,
		tryConnect: 0,
		tryChannel: 0,
	}
}

// AddHandler Adds a message handler to the consumer
func (c *Consumer) AddHandler(h func(amqp.Delivery)) {
	c.handlers = append(c.handlers, h)
}

// Start Calls doStart in a loop to handle connection loss
func (c *Consumer) Start() {

	for {
		c.doStart()

		c.log.Info("Waiting 1 min before retrying")
		time.Sleep(retryDelay)
	}
}

// Sets connection up and consume messages
func (c *Consumer) doStart() {

	defer c.tearDown()
	err := c.setup()
	if err != nil {
		c.log.Errorf("Error while setting up: %v", err)
		return
	}

	c.consume()
}

func (c *Consumer) setup() (err error) {

	// Connect to server
	c.tryConnect++
	//c.conn, err = amqp.Dial(c.amqpURL)
	c.conn, err = amqp.DialConfig(c.amqpURL, amqpConfig)
	if err != nil {
		return fmt.Errorf("Failed to connect to RabbitMQ: try %d, %v", c.tryConnect, err)
	}
	c.tryConnect = 0
	c.log.Infof("Consumer connected to AMQP server at: %s", c.amqpLogURL)

	// Open channel
	c.tryChannel++
	c.ch, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: try %d, %v", c.tryChannel, err)
	}
	c.tryChannel = 0

	c.deliveryChan, err = c.bindToQueue()
	if err != nil {
		c.log.Error(err)
		return err
	}

	return nil
}

func (c *Consumer) tearDown() {

	if c.ch != nil {
		c.ch.Close()
	}

	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) consume() {

	foreverChan := make(chan bool)

	// Loop to consume messages
	go c.doConsume()

	// Loop to check connection
	go checkConnection(c.log, c.conn, foreverChan, checkConnInterval)

	c.log.Info("Waiting for events . To exit press CTRL+C")
	<-foreverChan
}

func (c *Consumer) doConsume() {

	for d := range c.deliveryChan {
		for _, h := range c.handlers {
			h(d)
		}
	}
}

func (c *Consumer) bindToQueue() (<-chan amqp.Delivery, error) {

	q, err := declareQueue(c.ch, c.queueName)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare the queue '%s': %v", c.queueName, err)
	}
	c.log.Infof("Queue '%s' declared", c.queueName)

	var errs []error
	var targetCount int
	var boundCount int
	for exchangeName, routingKeys := range c.exchanges {
		for _, routingKey := range routingKeys {
			targetCount++

			err = bindQueue(c.ch, q, exchangeName, routingKey)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			boundCount++
			c.log.Infof("Bound to exchange '%s' with routing key '%s'", exchangeName, routingKey)
		}
	}

	if boundCount == 0 {
		return nil, fmt.Errorf("Couldn't bind routing keys. %d/%d'", boundCount, targetCount)
	}

	if boundCount < targetCount {
		c.log.Warnf("Failed to bind all routing keys %d/%d: %v", boundCount, targetCount, errs)
	}

	deliveryChan, err := consumeQueue(c.ch, q)
	if err != nil {
		return nil, fmt.Errorf("Failed to register the consumer: %v", err)
	}

	return deliveryChan, nil
}
