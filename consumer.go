package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/harik/weather_monitor/internal/api"
	"github.com/harik/weather_monitor/internal/evaluator"
	"github.com/harik/weather_monitor/internal/metrics"
)

func StartConsumer(brokers []string, topic string) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_0_0_0

	master, err := sarama.NewConsumerGroup(brokers, "weather-evaluator-group", config)
	if err != nil {
		return err
	}
	defer master.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	handler := &consumerGroupHandler{}

	for {
		if err := master.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("error from consumer: %v", err)
			time.Sleep(2 * time.Second)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
}

type consumerGroupHandler struct{}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		metrics.IncKafkaConsume()
		var data api.WeatherData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			log.Printf("invalid message: %v", err)
			continue
		}
		// call evaluator (non-blocking)
		go evaluator.Evaluate(&data)
		sess.MarkMessage(msg, "")
	}

	return nil
}
