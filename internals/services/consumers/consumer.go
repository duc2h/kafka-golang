package consumers

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

type Reader interface {
	Consume(ctx context.Context) error
}

func NewReader(config *sarama.Config, brokers, group, topics string, handler handler) Reader {
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Fatalf("Cannot create consumer group: %w", err)
	}

	return &reader{
		topics: strings.Split(topics, ","),
		consumerGroup: ConsumerGroup{
			onHandler: handler,
		},
		saramaConsumerGroup: client,
	}
}

type reader struct {
	topics              []string
	consumerGroup       ConsumerGroup
	saramaConsumerGroup sarama.ConsumerGroup
}

func (r *reader) Consume(ctx context.Context) error {
	return r.saramaConsumerGroup.Consume(ctx, r.topics, &r.consumerGroup)
}

// consumer-group

type handler func(*Message) error

type ConsumerGroup struct {
	onHandler handler
}

type Message struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64

	Timestamp time.Time
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			msg := Message{
				Key:       message.Key,
				Value:     message.Value,
				Topic:     message.Topic,
				Partition: message.Partition,
				Offset:    message.Offset,
				Timestamp: message.Timestamp,
			}
			err := c.onHandler(&msg)

			if err != nil {
				// add logic to control error.
				log.Printf("onHandler message error: %s", err.Error())
				continue
			}

			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
