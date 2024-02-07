package rabbitmq

import (
	"encoding/json"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

type Producer interface {
	SendMessage(key, contentType string, msg interface{}) error
}

type producer struct {
	qCfg        domain.AMPQStructCfg
	conn        Connection
	mu          sync.Mutex
	isConnected bool
}

var _ Producer = (*producer)(nil)

func NewProducer(cfg domain.AMPQStructCfg, conn Connection) *producer {
	return &producer{qCfg: cfg, conn: conn}
}

func (p *producer) upExchanges() error {
	if err := p.conn.ExchangeDeclare(
		p.qCfg.WorkExchange.Name,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := p.conn.ExchangeDeclare(
		p.qCfg.RetryExchange.Name,
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

func (p *producer) upQueues() error {

	for _, q := range p.qCfg.Queries {
		t := make(map[string]interface{})
		t["x-dead-letter-exchange"] = q.DLX
		if q.TTL > 0 {
			t["x-message-ttl"] = q.TTL
		}
		_, err := p.conn.QueueDeclare(
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

func (p *producer) upBinds() error {

	for _, b := range p.qCfg.QueryBinds {
		if err := p.conn.QueueBind(
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

func (p *producer) connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isConnected {
		return nil
	}

	if err := p.upExchanges(); err != nil {
		return err
	}

	if err := p.upQueues(); err != nil {
		return err
	}

	if err := p.upBinds(); err != nil {
		return err
	}
	p.isConnected = true
	return nil
}

func (p *producer) SendMessage(key, contentType string, msg interface{}) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if !p.isConnected {
		if err = p.connect(); err != nil {
			logrus.Errorf("Connect publisher to RabbitMQ error %s", err)
			return err
		}
	}

	if err = p.conn.Publish(
		p.qCfg.WorkExchange.Name,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        msgBytes,
		},
	); err != nil {
		p.mu.Lock()
		p.isConnected = false
		p.mu.Unlock()
		return err
	}
	return nil
}
