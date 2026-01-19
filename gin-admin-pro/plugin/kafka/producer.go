package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// Producer Kafka生产者
type Producer struct {
	client    *KafkaClient
	syncProd  sarama.SyncProducer
	asyncProd sarama.AsyncProducer

	config *ProducerConfig
	mu     sync.RWMutex
	ready  bool
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
	// 基础配置
	ClientID string   `yaml:"clientId" mapstructure:"clientId"`
	Brokers  []string `yaml:"brokers" mapstructure:"brokers"`

	// 发送配置
	RequiredAcks int16         `yaml:"requiredAcks" mapstructure:"requiredAcks"`
	Timeout      time.Duration `yaml:"timeout" mapstructure:"timeout"`
	Retries      int           `yaml:"retries" mapstructure:"retries"`
	RetryBackoff time.Duration `yaml:"retryBackoff" mapstructure:"retryBackoff"`

	// 压缩配置
	Compression string `yaml:"compression" mapstructure:"compression"`

	// 分区配置
	Partitioner string `yaml:"partitioner" mapstructure:"partitioner"`

	// 刷新配置
	FlushFrequency time.Duration `yaml:"flushFrequency" mapstructure:"flushFrequency"`
	FlushBytes     int           `yaml:"flushBytes" mapstructure:"flushBytes"`
	FlushMessages  int           `yaml:"flushMessages" mapstructure:"flushMessages"`

	// 消息配置
	MaxMessageBytes int `yaml:"maxMessageBytes" mapstructure:"maxMessageBytes"`

	// 返回配置
	ReturnSuccesses bool `yaml:"returnSuccesses" mapstructure:"returnSuccesses"`
	ReturnErrors    bool `yaml:"returnErrors" mapstructure:"returnErrors"`

	// 错误处理
	ErrorHandler func(error) `yaml:"-" mapstructure:"-"`

	// 成功处理
	SuccessHandler func(*sarama.ProducerMessage) `yaml:"-" mapstructure:"-"`
}

// NewProducer 创建生产者
func NewProducer(config *ProducerConfig) (*Producer, error) {
	if config == nil {
		config = DefaultProducerConfig()
	}

	producer := &Producer{
		config: config,
	}

	return producer, nil
}

// DefaultProducerConfig 返回默认生产者配置
func DefaultProducerConfig() *ProducerConfig {
	return &ProducerConfig{
		ClientID:        "gin-admin-producer",
		Brokers:         []string{"localhost:9092"},
		RequiredAcks:    1, // WaitForLocal
		Timeout:         time.Second * 30,
		Retries:         3,
		RetryBackoff:    time.Second * 2,
		Compression:     "none",
		Partitioner:     "hash",
		FlushFrequency:  time.Millisecond * 100,
		FlushBytes:      16384,
		FlushMessages:   0,
		MaxMessageBytes: 1000000, // 1MB
		ReturnSuccesses: true,
		ReturnErrors:    true,
	}
}

// Initialize 初始化生产者
func (p *Producer) Initialize() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 创建同步生产者
	syncConfig, err := p.buildSyncConfig()
	if err != nil {
		return fmt.Errorf("build sync config failed: %w", err)
	}

	p.syncProd, err = sarama.NewSyncProducer(p.config.Brokers, syncConfig)
	if err != nil {
		return fmt.Errorf("create sync producer failed: %w", err)
	}

	// 创建异步生产者
	asyncConfig, err := p.buildAsyncConfig()
	if err != nil {
		return fmt.Errorf("build async config failed: %w", err)
	}

	p.asyncProd, err = sarama.NewAsyncProducer(p.config.Brokers, asyncConfig)
	if err != nil {
		return fmt.Errorf("create async producer failed: %w", err)
	}

	// 启动异步处理
	go p.handleAsyncMessages()

	p.ready = true
	log.Printf("Producer initialized successfully")

	return nil
}

// buildSyncConfig 构建同步生产者配置
func (p *Producer) buildSyncConfig() (*sarama.Config, error) {
	return p.buildConfig()
}

// buildAsyncConfig 构建异步生产者配置
func (p *Producer) buildAsyncConfig() (*sarama.Config, error) {
	return p.buildConfig()
}

// buildConfig 构建生产者配置
func (p *Producer) buildConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()

	// 客户端配置
	config.ClientID = p.config.ClientID

	// 网络配置
	config.Net.DialTimeout = p.config.Timeout
	config.Net.ReadTimeout = p.config.Timeout
	config.Net.WriteTimeout = p.config.Timeout

	// 生产者配置
	config.Producer.RequiredAcks = sarama.RequiredAcks(p.config.RequiredAcks)
	config.Producer.Retry.Max = p.config.Retries
	config.Producer.Retry.Backoff = p.config.RetryBackoff
	config.Producer.Return.Successes = p.config.ReturnSuccesses
	config.Producer.Return.Errors = p.config.ReturnErrors

	// 压缩配置
	switch p.config.Compression {
	case "gzip":
		config.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		config.Producer.Compression = sarama.CompressionSnappy
	case "lz4":
		config.Producer.Compression = sarama.CompressionLZ4
	case "zstd":
		config.Producer.Compression = sarama.CompressionZSTD
	default:
		config.Producer.Compression = sarama.CompressionNone
	}

	// 分区器配置
	switch p.config.Partitioner {
	case "random":
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	case "roundrobin":
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case "hash":
		config.Producer.Partitioner = sarama.NewHashPartitioner
	default:
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}

	// 刷新配置
	config.Producer.Flush.Frequency = p.config.FlushFrequency
	config.Producer.Flush.Bytes = p.config.FlushBytes
	config.Producer.Flush.Messages = p.config.FlushMessages

	// 消息配置
	config.Producer.MaxMessageBytes = p.config.MaxMessageBytes

	return config, nil
}

// handleAsyncMessages 处理异步消息
func (p *Producer) handleAsyncMessages() {
	for {
		select {
		case errMsg := <-p.asyncProd.Errors():
			if p.config.ErrorHandler != nil {
				p.config.ErrorHandler(errMsg)
			} else {
				log.Printf("Producer error: %v", errMsg)
			}
		case successMsg := <-p.asyncProd.Successes():
			if p.config.SuccessHandler != nil {
				p.config.SuccessHandler(successMsg)
			}
		}
	}
}

// SendMessageSync 同步发送消息
func (p *Producer) SendMessageSync(ctx context.Context, topic string, key, value []byte) (int32, int64, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.ready {
		return 0, 0, fmt.Errorf("producer is not ready")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("timestamp"),
				Value: []byte(time.Now().Format(time.RFC3339)),
			},
		},
	}

	partition, offset, err := p.syncProd.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("send message failed: %w", err)
	}

	log.Printf("Message sent successfully to topic: %s, partition: %d, offset: %d", topic, partition, offset)
	return partition, offset, nil
}

// SendMessageAsync 异步发送消息
func (p *Producer) SendMessageAsync(ctx context.Context, topic string, key, value []byte) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.ready {
		return fmt.Errorf("producer is not ready")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("timestamp"),
				Value: []byte(time.Now().Format(time.RFC3339)),
			},
		},
	}

	select {
	case p.asyncProd.Input() <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// SendBatchMessagesSync 批量同步发送消息
func (p *Producer) SendBatchMessagesSync(ctx context.Context, topic string, messages []*Message) error {
	for _, msg := range messages {
		if _, _, err := p.SendMessageSync(ctx, topic, msg.Key, msg.Value); err != nil {
			return fmt.Errorf("send batch message failed: %w", err)
		}
	}
	return nil
}

// SendBatchMessagesAsync 批量异步发送消息
func (p *Producer) SendBatchMessagesAsync(ctx context.Context, topic string, messages []*Message) error {
	for _, msg := range messages {
		if err := p.SendMessageAsync(ctx, topic, msg.Key, msg.Value); err != nil {
			return fmt.Errorf("send batch message async failed: %w", err)
		}
	}
	return nil
}

// Close 关闭生产者
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.ready {
		return nil
	}

	p.ready = false

	var errs []error

	if p.syncProd != nil {
		if err := p.syncProd.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close sync producer failed: %w", err))
		}
	}

	if p.asyncProd != nil {
		if err := p.asyncProd.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close async producer failed: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close producer errors: %v", errs)
	}

	log.Println("Producer closed successfully")
	return nil
}

// IsReady 检查生产者是否就绪
func (p *Producer) IsReady() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ready
}

// GetConfig 获取配置
func (p *Producer) GetConfig() *ProducerConfig {
	return p.config
}
