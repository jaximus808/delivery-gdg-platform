package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Receive external events from other parts
//

func createTopicsIfNotExist(brokers string, topics []string) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return fmt.Errorf("failed to create admin client: %w", err)
	}
	defer adminClient.Close()

	// Prepare topic specifications
	var topicSpecs []kafka.TopicSpecification
	for _, topic := range topics {
		topicSpecs = append(topicSpecs, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := adminClient.CreateTopics(
		ctx,
		topicSpecs,
		kafka.SetAdminOperationTimeout(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	// Check results
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Printf("Failed to create topic %s: %v", result.Topic, result.Error)
		} else {
			log.Printf("Topic %s ready", result.Topic)
		}
	}

	return nil
}

func CreateKafkaConsumer(brokers, clientID string, topics []string) (*kafka.Consumer, error) {
	err := createTopicsIfNotExist(brokers, topics)
	if err != nil {
		log.Printf("Warning: failed to create topics: %v", err)
		// Continue anyway - topics might already exist or auto-create might be enabled
	}

	config := &kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"group.id":           clientID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	return consumer, nil
}
