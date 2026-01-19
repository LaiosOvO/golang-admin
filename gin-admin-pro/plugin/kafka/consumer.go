package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// Consumer Kafka消费者
type Consumer struct {
	client *KafkaClient
	group  sarama.ConsumerGroup

	config *ConsumerConfig
	mu     sync.RWMutex
	ready  bool

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	handlers map[string]sarama.ConsumerGroupHandler
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	// 基础配置
	ClientID string   `yaml:"clientId" mapstructure:"clientId"`
	Brokers  []string `yaml:"brokers" mapstructure:"brokers"`
	GroupID  string   `yaml:"groupId" mapstructure:"groupId"`

	// 消费配置
	AutoOffset string `yaml:"autoOffset" mapstructure:"autoOffset"` // earliest, latest, none
	AutoCommit bool   `yaml:"autoCommit" mapstructure:"autoCommit"`

	// 网络配置
	DialTimeout  time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout" mapstructure:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" mapstructure:"writeTimeout"`

	// 会话配置
	SessionTimeout    time.Duration `yaml:"sessionTimeout" mapstructure:"sessionTimeout"`
	HeartbeatInterval time.Duration `yaml:"heartbeatInterval" mapstructure:"heartbeatInterval"`
	RebalanceTimeout  time.Duration `yaml:"rebalanceTimeout" mapstructure:"rebalanceTimeout"`

	// 偏移量配置
	OffsetCommitInterval time.Duration `yaml:"offsetCommitInterval" mapstructure:"offsetCommitInterval"`

	// 获取配置
	FetchMin          int32         `yaml:"fetchMin" mapstructure:"fetchMin"`
	FetchDefault      int32         `yaml:"fetchDefault" mapstructure:"fetchDefault"`
	FetchMax          int32         `yaml:"fetchMax" mapstructure:"fetchMax"`
	MaxProcessingTime time.Duration `yaml:"maxProcessingTime" mapstructure:"maxProcessingTime"`
	MaxWaitTime       time.Duration `yaml:"maxWaitTime" mapstructure:"maxWaitTime"`

	// 错误处理
	ErrorHandler func(error) `yaml:"-" mapstructure:"-"`

	// 重新平衡处理
	RebalanceHandler func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error `yaml:"-" mapstructure:"-"`
}

// NewConsumer 创建消费者
func NewConsumer(config *ConsumerConfig) (*Consumer, error) {
	if config == nil {
		config = DefaultConsumerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := &Consumer{
		config:   config,
		ctx:      ctx,
		cancel:   cancel,
		handlers: make(map[string]sarama.ConsumerGroupHandler),
	}

	return consumer, nil
}

// DefaultConsumerConfig 返回默认消费者配置
func DefaultConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		ClientID:             "gin-admin-consumer",
		Brokers:              []string{"localhost:9092"},
		GroupID:              "gin-admin-group",
		AutoOffset:           "latest",
		AutoCommit:           true,
		DialTimeout:          time.Second * 30,
		ReadTimeout:          time.Second * 30,
		WriteTimeout:         time.Second * 30,
		SessionTimeout:       time.Second * 10,
		HeartbeatInterval:    time.Second * 3,
		RebalanceTimeout:     time.Second * 60,
		OffsetCommitInterval: time.Second,
		FetchMin:             1,
		FetchDefault:         1024 * 1024,      // 1MB
		FetchMax:             10 * 1024 * 1024, // 10MB
		MaxProcessingTime:    time.Second * 30,
		MaxWaitTime:          time.Millisecond * 500,
	}
}

// Initialize 初始化消费者
func (c *Consumer) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	config, err := c.buildConfig()
	if err != nil {
		return fmt.Errorf("build consumer config failed: %w", err)
	}

	c.group, err = sarama.NewConsumerGroup(c.config.Brokers, c.config.GroupID, config)
	if err != nil {
		return fmt.Errorf("create consumer group failed: %w", err)
	}

	c.ready = true
	log.Printf("Consumer initialized successfully with group: %s", c.config.GroupID)

	return nil
}

// buildConfig 构建消费者配置
func (c *Consumer) buildConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()

	// 客户端配置
	config.ClientID = c.config.ClientID

	// 网络配置
	config.Net.DialTimeout = c.config.DialTimeout
	config.Net.ReadTimeout = c.config.ReadTimeout
	config.Net.WriteTimeout = c.config.WriteTimeout

	// 消费者组配置
	config.Consumer.Group.Session.Timeout = c.config.SessionTimeout
	config.Consumer.Group.Heartbeat.Interval = c.config.HeartbeatInterval
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Group.ResetInvalidOffsets = true

	// 偏移量配置
	switch c.config.AutoOffset {
	case "earliest":
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "latest":
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	case "none":
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	// 提交配置
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = c.config.OffsetCommitInterval

	// 获取配置
	config.Consumer.Fetch.Min = c.config.FetchMin
	config.Consumer.Fetch.Default = c.config.FetchDefault
	config.Consumer.Fetch.Max = c.config.FetchMax
	config.Consumer.MaxProcessingTime = c.config.MaxProcessingTime
	config.Consumer.MaxWaitTime = c.config.MaxWaitTime

	return config, nil
}

// Subscribe 订阅主题
func (c *Consumer) Subscribe(topics []string, handler sarama.ConsumerGroupHandler) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.ready {
		return fmt.Errorf("consumer is not ready")
	}

	// 启动消费goroutine
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.runConsumer(topics, handler)
	}()

	log.Printf("Subscribed to topics: %v", topics)
	return nil
}

// runConsumer 运行消费者
func (c *Consumer) runConsumer(topics []string, handler sarama.ConsumerGroupHandler) {
	for {
		select {
		case <-c.ctx.Done():
			log.Printf("Consumer stopped")
			return
		default:
			if err := c.group.Consume(c.ctx, topics, handler); err != nil {
				if c.config.ErrorHandler != nil {
					c.config.ErrorHandler(err)
				} else {
					log.Printf("Error from consumer: %v", err)
				}
				time.Sleep(time.Second * 5) // 等待5秒后重试
			}
		}
	}
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.ready {
		return nil
	}

	c.ready = false

	// 取消上下文
	c.cancel()

	// 等待所有goroutine结束
	c.wg.Wait()

	// 关闭消费者组
	if c.group != nil {
		if err := c.group.Close(); err != nil {
			return fmt.Errorf("close consumer group failed: %w", err)
		}
	}

	log.Println("Consumer closed successfully")
	return nil
}

// IsReady 检查消费者是否就绪
func (c *Consumer) IsReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ready
}

// GetConfig 获取配置
func (c *Consumer) GetConfig() *ConsumerConfig {
	return c.config
}

// SimpleConsumerHandler 简单消费者处理器
type SimpleConsumerHandler struct {
	MessageHandler func(topic string, partition int32, offset int64, key, value []byte) error
}

// Setup Setup实现
func (h *SimpleConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup Cleanup实现
func (h *SimpleConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消息处理实现
func (h *SimpleConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.MessageHandler(message.Topic, message.Partition, message.Offset, message.Key, message.Value); err != nil {
			log.Printf("Process message failed: %v", err)
			return err
		}

		// 标记消息已处理
		session.MarkMessage(message, "")
	}
	return nil
}

// BatchConsumerHandler 批量消费者处理器
type BatchConsumerHandler struct {
	BatchSize      int
	BatchTimeout   time.Duration
	MessageHandler func(messages []*sarama.ConsumerMessage) error
}

// Setup Setup实现
func (h *BatchConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup Cleanup实现
func (h *BatchConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消息处理实现
func (h *BatchConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batch := make([]*sarama.ConsumerMessage, 0, h.BatchSize)
	timer := time.NewTimer(h.BatchTimeout)
	defer timer.Stop()

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				// 通道关闭，处理剩余的批量消息
				if len(batch) > 0 {
					if err := h.processBatch(session, batch); err != nil {
						return err
					}
				}
				return nil
			}

			batch = append(batch, message)

			// 达到批量大小，立即处理
			if len(batch) >= h.BatchSize {
				if err := h.processBatch(session, batch); err != nil {
					return err
				}
				batch = batch[:0] // 清空批量
				timer.Reset(h.BatchTimeout)
			}

		case <-timer.C:
			// 超时，处理当前批量
			if len(batch) > 0 {
				if err := h.processBatch(session, batch); err != nil {
					return err
				}
				batch = batch[:0] // 清空批量
			}
			timer.Reset(h.BatchTimeout)
		}
	}
}

// processBatch 处理批量消息
func (h *BatchConsumerHandler) processBatch(session sarama.ConsumerGroupSession, batch []*sarama.ConsumerMessage) error {
	if err := h.MessageHandler(batch); err != nil {
		return err
	}

	// 标记所有消息已处理
	for _, message := range batch {
		session.MarkMessage(message, "")
	}

	return nil
}
