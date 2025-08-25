package consumer

import (
	"context"
	"log"

	"github.com/Sphirium/learning-projects/wb-tech-demo-lo/internal/service"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	service *service.OrderService
}

func NewKafkaConsumer(broker, topic string, service *service.OrderService) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  "order-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		reader:  reader,
		service: service,
	}
}

func (c *KafkaConsumer) Start() {
	go func() {
		for {
			msg, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			log.Printf("📨 Получено сообщение: key=%s, value=%s", string(msg.Key), string(msg.Value))

			if err := c.service.SaveOrder(msg.Value); err != nil {
				log.Printf("❌ Failed to process order: %v", err)
				continue
			}

			log.Printf("✅ Успешно обработан заказ: %s", msg.Key)
		}
	}()
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
