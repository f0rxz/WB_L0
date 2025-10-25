package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewConsumer(reader *kafka.Reader, writer *kafka.Writer) *Consumer {
	return &Consumer{reader: reader, writer: writer}
}

func getRetryCount(msg kafka.Message) int {
	for _, h := range msg.Headers {
		if h.Key == "x-retry-count" {
			var count int
			fmt.Sscanf(string(h.Value), "%d", &count)
			return count
		}
	}
	return 0
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
			retryCount := getRetryCount(msg) + 1

			if retryCount <= 3 {
				log.Printf("kafka: handler failed, retry #%d for offset %d", retryCount, msg.Offset)

				retryMsg := kafka.Message{
					Key:   msg.Key,
					Value: msg.Value,
					Headers: []kafka.Header{
						{Key: "x-retry-count", Value: []byte(fmt.Sprintf("%d", retryCount))},
					},
				}

				if err := c.writer.WriteMessages(ctx, retryMsg); err != nil {
					log.Printf("kafka: failed to publish retry message: %v", err)
				}
			} else {
				log.Printf("kafka: sending message offset %d to DLQ", msg.Offset)
				dlqMsg := kafka.Message{
					Key:   msg.Key,
					Value: msg.Value,
					Headers: []kafka.Header{
						{Key: "x-original-topic", Value: []byte(c.reader.Config().Topic)},
					},
				}

				if err := c.writer.WriteMessages(ctx, dlqMsg); err != nil {
					log.Printf("kafka: failed to send to DLQ: %v", err)
				}
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("kafka: failed to commit after error: %v", err)
			}

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
