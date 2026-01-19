package kafka

import "github.com/IBM/sarama"

// Message 消息结构
type Message struct {
	Key     []byte
	Value   []byte
	Headers map[string]string
}

// ConsumerHandler 消费者处理器接口
type ConsumerHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}
