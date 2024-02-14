package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VadimGossip/calculator/agent/internal/domain"
	"github.com/VadimGossip/calculator/agent/internal/worker"
	"github.com/VadimGossip/calculator/agent/pkg/workerctrl"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Consumer interface {
	Subscribe(ctx context.Context)
}

type consumer struct {
	cfg           domain.AMPQStructCfg
	conn          Connection
	workerService worker.Service
}

var _ Consumer = (*consumer)(nil)

func NewConsumer(cfg domain.AMPQStructCfg, conn Connection, workerService worker.Service) *consumer {
	return &consumer{cfg: cfg, conn: conn, workerService: workerService}
}

func (c *consumer) upExchanges() error {
	if err := c.conn.ExchangeDeclare(
		c.cfg.WorkExchange.Name,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := c.conn.ExchangeDeclare(
		c.cfg.RetryExchange.Name,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (c *consumer) upQueues() error {

	for _, q := range c.cfg.Queries {
		t := make(map[string]interface{})
		t["x-dead-letter-exchange"] = q.DLX
		if q.TTL > 0 {
			t["x-message-ttl"] = q.TTL
		}
		_, err := c.conn.QueueDeclare(
			q.Name,
			true,
			false,
			false,
			false,
			t,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *consumer) upBinds() error {

	for _, b := range c.cfg.QueryBinds {
		if err := c.conn.QueueBind(
			b.QueryName,
			b.Key,
			b.ExchangeName,
			false,
			nil,
		); err != nil {
			return err
		}
	}
	return nil
}

func (c *consumer) connect() (<-chan amqp.Delivery, error) {
	if err := c.upExchanges(); err != nil {
		return nil, err
	}

	if err := c.upQueues(); err != nil {
		return nil, err
	}

	if err := c.upBinds(); err != nil {
		return nil, err
	}

	msg, err := c.conn.Consume(
		c.cfg.ConsumerCfg.QueryName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *consumer) subscribe(ctx context.Context) {
	var msg <-chan amqp.Delivery
	var err error

	for {
		if msg, err = c.connect(); err != nil {
			logrus.Errorf("Connect consumer to RabbitMQ error %s", err)
			time.Sleep(15 * time.Second)
			continue
		}
		break
	}
	logrus.Info("RabbitMQ consumer connected")
	workerCtrl := workerctrl.NewService(c.workerService.GetMaxProcessAllowed())
	for {
		select {
		case <-ctx.Done():
			logrus.Info("RabbitMQ consumer stopped")
			return
		case v, ok := <-msg:
			if !ok {
				logrus.Info("Consumer msg channel closed")
				if c.conn.IsClosed() {
					return
				}
				logrus.Info("Try reconnect")
				go func() {
					c.subscribe(ctx)
				}()
				return
			}
			workerCtrl.Acquire(1)
			go func(v amqp.Delivery) {
				defer workerCtrl.Release(1)
				c.processDeliveryMsg(v)
			}(v)
			var id int64
			if err = json.Unmarshal(v.Body, &id); err != nil {
				fmt.Println(err)
			}

			fmt.Println(id)
		}
	}
}

func (c *consumer) processDeliveryMsg(msg amqp.Delivery) {
	var item domain.SubExpressionQueryItem
	if err := json.Unmarshal(msg.Body, &item); err != nil {
		logrus.Errorf("unmarshal msg error. msg body %s error %s. Message will be thrown away", msg.Body, err)
	}
	if err := c.workerService.Do(item); err != nil {
		logrus.Errorf("unmarshal msg error. msg body %s error %s. Message will be moved to dlx", msg.Body, err)
		if err = msg.Nack(false, false); err != nil {
			logrus.Errorf("nack msg error %s", err)
		}
	}
	if err := msg.Ack(false); err != nil {
		logrus.Errorf("ack msg error %s", err)
	}
}

func (c *consumer) Subscribe(ctx context.Context) {
	go func() {
		c.subscribe(ctx)
	}()
}
