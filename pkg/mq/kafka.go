package mq

import (
	"context"
	"encoding/json"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	writer *kafka.Writer
}

func NewKafka() *Kafka {
	writer := &kafka.Writer{
		Addr:        kafka.TCP("localhost:9092"),
		Balancer:    &kafka.LeastBytes{},
		MaxAttempts: 3,
	}

	return &Kafka{
		writer: writer,
	}
}

func (k *Kafka) ProduceMessage(topic string, value interface{}) error {
	bvalue, err := json.Marshal(value)

	err = k.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Topic: topic,
			Value: bvalue,
		},
	)
	if err != nil {
		logger.Errorf("failed to write message: %s", err)
		return err
	}

	return nil
}

func (k *Kafka) StartConsumeMessages(topic, groupID string, f func(data []byte) error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10,   // 最小批量拉取
		MaxBytes: 10e6, // 最大批量拉取
	})
	defer reader.Close()
	for {
		m, err := reader.FetchMessage(context.Background())
		if err != nil {
			logger.Errorf("ReadMessage Fail: %s", err)
			continue
		}

		if err = f(m.Value); err != nil {
			logger.Errorf("ExecFunc Fail: %s", err)
			continue
		}

		// 处理成功 → 提交 offset
		if err = reader.CommitMessages(context.Background(), m); err != nil {
			logger.Errorf("commit failed: %v", err)
		}
	}
}
