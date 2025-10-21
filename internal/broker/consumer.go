package broker

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(reader *kafka.Reader) *Consumer {
	return &Consumer{reader: reader}
}

func (c *Consumer) Consume(ctx context.Context, handler func(ctx context.Context, key, value []byte) error) error {
	if c.reader == nil {
		return fmt.Errorf("kafka: reader is nil")
	}

	log.Println("kafka: consumer started")
	for {
		select {
		case <-ctx.Done():
			log.Println("kafka: consumer stopping due to context done")
			return nil
		default:
		}

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("kafka: context canceled, stopping consumer loop")
				return nil
			}
			log.Printf("kafka: fetch error: %v\n", err)
			continue
		}

		if err := handler(ctx, msg.Key, msg.Value); err != nil {
			log.Printf("kafka: handler failed, message offset %d will be retried: %v\n", msg.Offset, err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("kafka: failed to commit offset %d: %v\n", msg.Offset, err)
		}
	}
}

func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
