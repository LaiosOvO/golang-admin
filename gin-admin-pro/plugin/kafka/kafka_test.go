package kafka

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, []string{"localhost:9092"}, config.Brokers)
	assert.Equal(t, "gin-admin-group", config.GroupID)
	assert.True(t, config.AutoCommit)
	assert.Equal(t, "latest", config.AutoOffset)
}

func TestKafkaClient_NewKafkaClient(t *testing.T) {
	// Test with nil config
	client, err := NewKafkaClient(nil)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, []string{"localhost:9092"}, client.config.Brokers)

	// Test with custom config
	config := &Config{
		Brokers: []string{"broker1:9092", "broker2:9092"},
		GroupID: "test-group",
	}

	client, err = NewKafkaClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, []string{"broker1:9092", "broker2:9092"}, client.config.Brokers)
	assert.Equal(t, "test-group", client.config.GroupID)
}

func TestKafkaClient_BuildProducerConfig(t *testing.T) {
	client, _ := NewKafkaClient(nil)

	config, err := client.buildProducerConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "gin-admin-producer", config.ClientID)
	assert.Equal(t, sarama.RequiredAcks(1), config.Producer.RequiredAcks)
	assert.Equal(t, sarama.CompressionNone, config.Producer.Compression)
}

func TestKafkaClient_BuildConsumerConfig(t *testing.T) {
	client, _ := NewKafkaClient(nil)

	config, err := client.buildConsumerConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "gin-admin-consumer", config.ClientID)
	assert.Equal(t, sarama.OffsetNewest, config.Consumer.Offsets.Initial)
}

func TestProducer_NewProducer(t *testing.T) {
	// Test with nil config
	producer, err := NewProducer(nil)
	assert.NoError(t, err)
	assert.NotNil(t, producer)
	assert.Equal(t, "gin-admin-producer", producer.config.ClientID)

	// Test with custom config
	config := &ProducerConfig{
		ClientID: "test-producer",
		Brokers:  []string{"broker1:9092"},
	}

	producer, err = NewProducer(config)
	assert.NoError(t, err)
	assert.NotNil(t, producer)
	assert.Equal(t, "test-producer", producer.config.ClientID)
	assert.Equal(t, []string{"broker1:9092"}, producer.config.Brokers)
}

func TestProducer_BuildConfig(t *testing.T) {
	producer, _ := NewProducer(nil)

	config, err := producer.buildConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "gin-admin-producer", config.ClientID)
	assert.Equal(t, sarama.RequiredAcks(1), config.Producer.RequiredAcks)
}

func TestConsumer_NewConsumer(t *testing.T) {
	// Test with nil config
	consumer, err := NewConsumer(nil)
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, "gin-admin-consumer", consumer.config.ClientID)

	// Test with custom config
	config := &ConsumerConfig{
		ClientID: "test-consumer",
		GroupID:  "test-group",
		Brokers:  []string{"broker1:9092"},
	}

	consumer, err = NewConsumer(config)
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, "test-consumer", consumer.config.ClientID)
	assert.Equal(t, "test-group", consumer.config.GroupID)
}

func TestConsumer_BuildConfig(t *testing.T) {
	consumer, _ := NewConsumer(nil)

	config, err := consumer.buildConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "gin-admin-consumer", config.ClientID)
	assert.Equal(t, sarama.OffsetNewest, config.Consumer.Offsets.Initial)
}

func TestSimpleConsumerHandler(t *testing.T) {
	var processedMessages []string

	handler := &SimpleConsumerHandler{
		MessageHandler: func(topic string, partition int32, offset int64, key, value []byte) error {
			processedMessages = append(processedMessages, string(value))
			return nil
		},
	}

	// Test Setup
	err := handler.Setup(nil)
	assert.NoError(t, err)

	// Test Cleanup
	err = handler.Cleanup(nil)
	assert.NoError(t, err)

	// Test with mock messages would require more complex setup
	assert.NotNil(t, handler)
}

func TestBatchConsumerHandler(t *testing.T) {
	var processedBatches [][]*sarama.ConsumerMessage

	handler := &BatchConsumerHandler{
		BatchSize:    3,
		BatchTimeout: time.Second * 5,
		MessageHandler: func(messages []*sarama.ConsumerMessage) error {
			processedBatches = append(processedBatches, messages)
			return nil
		},
	}

	// Test Setup
	err := handler.Setup(nil)
	assert.NoError(t, err)

	// Test Cleanup
	err = handler.Cleanup(nil)
	assert.NoError(t, err)

	assert.Equal(t, 3, handler.BatchSize)
	assert.Equal(t, time.Second*5, handler.BatchTimeout)
}

func TestMessage(t *testing.T) {
	msg := &Message{
		Key:     []byte("test-key"),
		Value:   []byte("test-value"),
		Headers: map[string]string{"header1": "value1"},
	}

	assert.Equal(t, []byte("test-key"), msg.Key)
	assert.Equal(t, []byte("test-value"), msg.Value)
	assert.Equal(t, "value1", msg.Headers["header1"])
}

func TestDefaultProducerConfig(t *testing.T) {
	config := DefaultProducerConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "gin-admin-producer", config.ClientID)
	assert.Equal(t, []string{"localhost:9092"}, config.Brokers)
	assert.Equal(t, int16(1), config.RequiredAcks)
	assert.Equal(t, "none", config.Compression)
	assert.Equal(t, "hash", config.Partitioner)
	assert.True(t, config.ReturnSuccesses)
	assert.True(t, config.ReturnErrors)
}

func TestDefaultConsumerConfig(t *testing.T) {
	config := DefaultConsumerConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "gin-admin-consumer", config.ClientID)
	assert.Equal(t, []string{"localhost:9092"}, config.Brokers)
	assert.Equal(t, "gin-admin-group", config.GroupID)
	assert.Equal(t, "latest", config.AutoOffset)
	assert.True(t, config.AutoCommit)
	assert.Equal(t, time.Second*10, config.SessionTimeout)
	assert.Equal(t, time.Second*3, config.HeartbeatInterval)
}

// Integration tests (would require actual Kafka instance)
// These are skipped by default to avoid requiring external dependencies

func TestKafkaClient_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test would require a running Kafka instance
	// Uncomment and modify when you have Kafka available for testing
	/*
		config := &Config{
			Brokers: []string{"localhost:9092"},
			GroupID: "test-group",
		}

		client, err := NewKafkaClient(config)
		require.NoError(t, err)

		err = client.Initialize()
		require.NoError(t, err)
		defer client.Close()

		assert.True(t, client.IsReady())
	*/
}

func TestProducer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test would require a running Kafka instance
	// Uncomment and modify when you have Kafka available for testing
	/*
		config := &ProducerConfig{
			Brokers: []string{"localhost:9092"},
			ClientID: "test-producer",
		}

		producer, err := NewProducer(config)
		require.NoError(t, err)

		err = producer.Initialize()
		require.NoError(t, err)
		defer producer.Close()

		assert.True(t, producer.IsReady())

		// Test send message
		partition, offset, err := producer.SendMessageSync(context.Background(), "test-topic", []byte("key"), []byte("value"))
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, partition, int32(0))
		assert.GreaterOrEqual(t, offset, int64(0))
	*/
}
