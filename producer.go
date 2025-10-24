package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/harik/weather_monitor/internal/api"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 5

	p, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	producer := &Producer{asyncProducer: p, topic: topic}

	// handle errors (non-blocking)
	go func() {
		for err := range p.Errors() {
			log.Printf("producer error: %v", err)
		}
	}()

	return producer, nil
}

func (p *Producer) Publish(data *api.WeatherData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(b),
	}
	p.asyncProducer.Input() <- msg
	return nil
}

func (p *Producer) Close() error {
	return p.asyncProducer.Close()
}
