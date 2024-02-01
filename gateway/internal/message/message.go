package message

import (
	"context"
	"github.com/streadway/amqp"
)

type Producer interface {
	Publish(topic string, message string) error
	
	PublishWithContext(ctx context.Context, topic string, message string) error
}

type RabbitmqProducer struct {
	client  *amqp.Connection
	channel *amqp.Channel
}

func (r *RabbitmqProducer) Publish(topic string, message string) error {
	//TODO implement me
	panic("implement me")
}

func (r *RabbitmqProducer) PublishWithContext(ctx context.Context, topic string, message string) error {
	//TODO implement me
	panic("implement me")
}

func NewRabbitmqProducer(amqpUrl string) (Producer, error) {
	dial, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, err
	}
	channel, err := dial.Channel()
	if err != nil {
		return nil, err
	}
	r := &RabbitmqProducer{
		client:  dial,
		channel: channel,
	}
	return r, nil
}
