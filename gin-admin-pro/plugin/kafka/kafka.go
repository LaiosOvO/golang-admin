package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

var (
	ErrProducerNotReady = errors.New("producer is not ready")
	ErrConsumerNotReady = errors.New("consumer is not ready")
)

// KafkaClient Kafka客户端包装器
type KafkaClient struct {
	config   *Config
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	admin    sarama.ClusterAdmin

	// 消费者管理
	consumers map[string]*ConsumerGroupManager
	mu        sync.RWMutex

	// 状态
	ready     bool
	closeChan chan struct{}
	wg        sync.WaitGroup
}

// ConsumerGroupManager 消费者组管理器
type ConsumerGroupManager struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler sarama.ConsumerGroupHandler
	cancel  context.CancelFunc
	ctx     context.Context
}

// NewKafkaClient 创建Kafka客户端
func NewKafkaClient(cfg *Config) (*KafkaClient, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	client := &KafkaClient{
		config:    cfg,
		consumers: make(map[string]*ConsumerGroupManager),
		closeChan: make(chan struct{}),
	}

	return client, nil
}

// Initialize 初始化Kafka客户端
func (kc *KafkaClient) Initialize() error {
	// 创建生产者
	producerConfig, err := kc.buildProducerConfig()
	if err != nil {
		return fmt.Errorf("build producer config failed: %w", err)
	}

	kc.producer, err = sarama.NewSyncProducer(kc.config.Brokers, producerConfig)
	if err != nil {
		return fmt.Errorf("create producer failed: %w", err)
	}

	// 创建集群管理员
	adminConfig, err := kc.buildAdminConfig()
	if err != nil {
		return fmt.Errorf("build admin config failed: %w", err)
	}

	kc.admin, err = sarama.NewClusterAdmin(kc.config.Brokers, adminConfig)
	if err != nil {
		return fmt.Errorf("create cluster admin failed: %w", err)
	}

	kc.ready = true
	log.Printf("Kafka client initialized successfully with brokers: %v", kc.config.Brokers)

	return nil
}

// buildProducerConfig 构建生产者配置
func (kc *KafkaClient) buildProducerConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()

	// 基础配置
	config.ClientID = "gin-admin-producer"
	config.Net.DialTimeout = kc.config.DialTimeout
	config.Net.ReadTimeout = kc.config.ReadTimeout
	config.Net.WriteTimeout = kc.config.WriteTimeout

	// 生产者配置
	config.Producer.RequiredAcks = sarama.RequiredAcks(kc.config.Producer.RequiredAcks)
	config.Producer.Retry.Max = kc.config.MaxRetries
	config.Producer.Retry.Backoff = kc.config.RetryBackoff
	config.Producer.Return.Successes = kc.config.Producer.ReturnSuccesses
	config.Producer.Return.Errors = kc.config.Producer.ReturnErrors

	// 压缩配置
	switch kc.config.Producer.Compression {
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
	switch kc.config.Producer.Partitioner {
	case "random":
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	case "roundrobin":
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case "hash":
		config.Producer.Partitioner = sarama.NewHashPartitioner
	default:
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}

	// 其他配置
	config.Producer.Flush.Frequency = kc.config.Producer.FlushFrequency
	config.Producer.Flush.Bytes = kc.config.Producer.FlushBytes
	config.Producer.Flush.Messages = kc.config.Producer.FlushMessages
	config.Producer.MaxMessageBytes = kc.config.Producer.MaxMessageBytes

	// 安全配置
	if err := kc.configureSecurity(config); err != nil {
		return nil, err
	}

	return config, nil
}

// buildConsumerConfig 构建消费者配置
func (kc *KafkaClient) buildConsumerConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()

	// 基础配置
	config.ClientID = "gin-admin-consumer"
	config.Net.DialTimeout = kc.config.DialTimeout
	config.Net.ReadTimeout = kc.config.ReadTimeout
	config.Net.WriteTimeout = kc.config.WriteTimeout

	// 消费者组配置
	config.Consumer.Group.Session.Timeout = kc.config.Consumer.SessionTimeout
	config.Consumer.Group.Heartbeat.Interval = kc.config.Consumer.HeartbeatInterval
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Group.ResetInvalidOffsets = true

	// 偏移量配置
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	switch kc.config.AutoOffset {
	case "earliest":
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "latest":
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	case "none":
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	// 提交配置
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = kc.config.Consumer.OffsetCommitInterval

	// 获取配置
	config.Consumer.Fetch.Min = kc.config.Consumer.FetchMin
	config.Consumer.Fetch.Default = kc.config.Consumer.FetchDefault
	config.Consumer.Fetch.Max = kc.config.Consumer.FetchMax
	config.Consumer.MaxProcessingTime = kc.config.Consumer.MaxProcessingTime
	config.Consumer.MaxWaitTime = kc.config.Consumer.MaxWaitTime

	// 安全配置
	if err := kc.configureSecurity(config); err != nil {
		return nil, err
	}

	return config, nil
}

// buildAdminConfig 构建管理员配置
func (kc *KafkaClient) buildAdminConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()

	// 基础配置
	config.ClientID = "gin-admin-admin"
	config.Net.DialTimeout = kc.config.DialTimeout
	config.Net.ReadTimeout = kc.config.ReadTimeout
	config.Net.WriteTimeout = kc.config.WriteTimeout

	// 安全配置
	if err := kc.configureSecurity(config); err != nil {
		return nil, err
	}

	return config, nil
}

// configureSecurity 配置安全认证
func (kc *KafkaClient) configureSecurity(config *sarama.Config) error {
	if !kc.config.Security.Enabled {
		return nil
	}

	// SASL配置
	if kc.config.Security.Username != "" && kc.config.Security.Password != "" {
		switch kc.config.Security.Mechanism {
		case "SCRAM-SHA-256":
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case "SCRAM-SHA-512":
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		default:
			config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}

		config.Net.SASL.Enable = true
		config.Net.SASL.User = kc.config.Security.Username
		config.Net.SASL.Password = kc.config.Security.Password
	}

	// TLS配置
	if kc.config.Security.TLSEnabled {
		config.Net.TLS.Enable = true
		// 这里可以添加更详细的TLS配置
	}

	return nil
}

// SendMessage 发送消息
func (kc *KafkaClient) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	if !kc.ready {
		return ErrProducerNotReady
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

	partition, offset, err := kc.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("send message failed: %w", err)
	}

	log.Printf("Message sent successfully to topic: %s, partition: %d, offset: %d", topic, partition, offset)
	return nil
}

// SendBatchMessages 批量发送消息
func (kc *KafkaClient) SendBatchMessages(ctx context.Context, topic string, messages []*Message) error {
	if !kc.ready {
		return ErrProducerNotReady
	}

	for _, msg := range messages {
		if err := kc.SendMessage(ctx, topic, msg.Key, msg.Value); err != nil {
			return fmt.Errorf("send batch message failed: %w", err)
		}
	}

	return nil
}

// CreateConsumer 创建消费者组
func (kc *KafkaClient) CreateConsumer(groupID string, topics []string, handler sarama.ConsumerGroupHandler) error {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	if _, exists := kc.consumers[groupID]; exists {
		return fmt.Errorf("consumer group %s already exists", groupID)
	}

	consumerConfig, err := kc.buildConsumerConfig()
	if err != nil {
		return fmt.Errorf("build consumer config failed: %w", err)
	}

	consumer, err := sarama.NewConsumerGroup(kc.config.Brokers, groupID, consumerConfig)
	if err != nil {
		return fmt.Errorf("create consumer group failed: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	manager := &ConsumerGroupManager{
		group:   consumer,
		topics:  topics,
		handler: handler,
		cancel:  cancel,
		ctx:     ctx,
	}

	kc.consumers[groupID] = manager

	// 启动消费者
	kc.wg.Add(1)
	go func() {
		defer kc.wg.Done()
		kc.runConsumer(manager)
	}()

	log.Printf("Consumer group %s created for topics: %v", groupID, topics)
	return nil
}

// runConsumer 运行消费者
func (kc *KafkaClient) runConsumer(manager *ConsumerGroupManager) {
	for {
		select {
		case <-manager.ctx.Done():
			log.Printf("Consumer group stopped")
			return
		default:
			if err := manager.group.Consume(manager.ctx, manager.topics, manager.handler); err != nil {
				log.Printf("Error from consumer: %v", err)
				time.Sleep(time.Second * 5) // 等待5秒后重试
			}
		}
	}
}

// CloseConsumer 关闭消费者组
func (kc *KafkaClient) CloseConsumer(groupID string) error {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	manager, exists := kc.consumers[groupID]
	if !exists {
		return fmt.Errorf("consumer group %s not found", groupID)
	}

	manager.cancel()
	if err := manager.group.Close(); err != nil {
		return fmt.Errorf("close consumer group failed: %w", err)
	}

	delete(kc.consumers, groupID)
	log.Printf("Consumer group %s closed", groupID)

	return nil
}

// CreateTopic 创建主题
func (kc *KafkaClient) CreateTopic(topic string, partitions int32, replicationFactor int16) error {
	if kc.admin == nil {
		return errors.New("cluster admin is not initialized")
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries: map[string]*string{
			"cleanup.policy": stringPtr("delete"),
			"retention.ms":   stringPtr("604800000"), // 7 days
		},
	}

	err := kc.admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		return fmt.Errorf("create topic failed: %w", err)
	}

	log.Printf("Topic %s created with %d partitions and replication factor %d", topic, partitions, replicationFactor)
	return nil
}

// DeleteTopic 删除主题
func (kc *KafkaClient) DeleteTopic(topic string) error {
	if kc.admin == nil {
		return errors.New("cluster admin is not initialized")
	}

	err := kc.admin.DeleteTopic(topic)
	if err != nil {
		return fmt.Errorf("delete topic failed: %w", err)
	}

	log.Printf("Topic %s deleted", topic)
	return nil
}

// ListTopics 列出所有主题
func (kc *KafkaClient) ListTopics() (map[string]sarama.TopicDetail, error) {
	if kc.admin == nil {
		return nil, errors.New("cluster admin is not initialized")
	}

	metadata, err := kc.admin.ListTopics()
	if err != nil {
		return nil, fmt.Errorf("list topics failed: %w", err)
	}

	return metadata, nil
}

// Close 关闭Kafka客户端
func (kc *KafkaClient) Close() error {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	if !kc.ready {
		return nil
	}

	kc.ready = false

	// 关闭所有消费者
	for groupID := range kc.consumers {
		if err := kc.CloseConsumer(groupID); err != nil {
			log.Printf("Error closing consumer %s: %v", groupID, err)
		}
	}

	// 等待所有goroutine结束
	close(kc.closeChan)
	kc.wg.Wait()

	// 关闭生产者
	if kc.producer != nil {
		if err := kc.producer.Close(); err != nil {
			return fmt.Errorf("close producer failed: %w", err)
		}
	}

	// 关闭管理员
	if kc.admin != nil {
		if err := kc.admin.Close(); err != nil {
			return fmt.Errorf("close cluster admin failed: %w", err)
		}
	}

	log.Println("Kafka client closed successfully")
	return nil
}

// IsReady 检查客户端是否就绪
func (kc *KafkaClient) IsReady() bool {
	return kc.ready
}

// GetConfig 获取配置
func (kc *KafkaClient) GetConfig() *Config {
	return kc.config
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
