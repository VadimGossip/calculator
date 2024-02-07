package rabbitmq

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type Connection interface {
	Connect(ctx context.Context) error
	Close() error
	Connection() *amqp.Connection
	Channel() (*amqp.Channel, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	GetChannelFromPool(exchange, key, queue, consumer string) (*amqp.Channel, error)
	IsClosed() bool
}

type ChannelPoolItemKey struct {
	Queue    string
	Consumer string
	Exchange string
	Key      string
}

type connection struct {
	url            string
	mu             sync.RWMutex
	conn           *amqp.Connection
	chMu           sync.RWMutex
	serviceChannel *amqp.Channel
	channelPool    map[ChannelPoolItemKey]*amqp.Channel
	closed         bool
	chanCtx        context.Context
}

var _ Connection = (*connection)(nil)

func NewConnection(url string) *connection {
	return &connection{url: url}
}

func (c *connection) Connection() *amqp.Connection {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

func (c *connection) Channel() (*amqp.Channel, error) {
	c.chMu.RLock()
	defer c.chMu.RUnlock()

	channel, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (c *connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.closed = true
	for _, ch := range c.channelPool {
		if err := ch.Close(); err != nil {
			return err
		}
	}

	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (c *connection) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

func (c *connection) connect() error {
	connectRabbitMQ, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	c.conn = connectRabbitMQ

	c.serviceChannel, err = c.conn.Channel()
	if err != nil {
		return err
	}
	c.channelPool = make(map[ChannelPoolItemKey]*amqp.Channel)
	return nil
}

func (c *connection) Connect(ctx context.Context) error {
	if !c.IsClosed() {
		if err := c.connect(); err != nil {
			return err
		}
	}
	c.chanCtx = ctx
	go func() {
		logrus.Info("Start watching RabbitMQ connection")
		for {
			select {
			case <-ctx.Done():
				logrus.Info("Stop watching RabbitMQ connection")
				return
			case _, ok := <-c.conn.NotifyClose(make(chan *amqp.Error)):
				if !ok {
					if c.IsClosed() {
						return
					}
					logrus.Infof("Unexpected close of RabbitMQ connection")

					c.mu.Lock()
					for {
						if err := c.connect(); err != nil {
							logrus.Info("connection failed, trying to reconnect to RabbitMQ")
							time.Sleep(time.Second * 15)
							continue
						}
						logrus.Info("RabbitMQ connection restored")
						break
					}
					c.mu.Unlock()
				}
			}
		}
	}()
	return nil
}

func (c *connection) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.serviceChannel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (c *connection) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.serviceChannel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (c *connection) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.serviceChannel.QueueBind(name, key, exchange, noWait, args)
}

func (c *connection) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, err := c.GetChannelFromPool("", "", queue, consumer)
	if err != nil {
		return nil, err
	}

	return ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (c *connection) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, err := c.GetChannelFromPool(exchange, key, "", "")
	if err != nil {
		return err
	}

	return ch.Publish(exchange, key, mandatory, immediate, msg)
}

func (c *connection) GetChannelFromPool(exchange, key, queue, consumer string) (*amqp.Channel, error) {
	c.chMu.Lock()
	defer c.chMu.Unlock()
	var err error
	poolKey := ChannelPoolItemKey{
		Exchange: exchange,
		Key:      key,
		Queue:    queue,
		Consumer: consumer,
	}
	channel, ok := c.channelPool[poolKey]
	if !ok {
		channel, err = c.conn.Channel()
		if err != nil {
			return nil, err
		}
		c.channelPool[poolKey] = channel
		go func() {
			logrus.Info("Start watching RabbitMQ channel")
			for {
				select {
				case <-c.chanCtx.Done():
					logrus.Info("Stop watching RabbitMQ channel")
					return
				case _, ok := <-channel.NotifyClose(make(chan *amqp.Error)):
					if !ok {
						if c.IsClosed() {
							return
						}
						logrus.Infof("Unexpected close of RabbitMQ channel")
						c.chMu.Lock()
						delete(c.channelPool, poolKey)
						c.chMu.Unlock()
						return
					}
				}
			}
		}()
	}

	return channel, nil
}
