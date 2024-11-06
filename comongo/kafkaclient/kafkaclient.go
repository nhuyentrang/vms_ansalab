package kafkaclient

import (
	"errors"
	"fmt"
	"strings"
	"time"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// consumerSessionTimeoutms represents the session timeout in milliseconds for the Kafka consumer.
var consumerSessionTimeoutms int = 60000

// kafkaJsonMessage represents the structure of a JSON message to be sent to Kafka.
type kafkaJsonMessage struct {
	Message   string
	Topic     string
	Timestamp time.Time
}

// _consumerList is a map that stores the Kafka consumers.
var _consumerList = make(map[string]*confluentKafka.Consumer)

// _publishMsgChannel is a channel used for publishing messages to Kafka.
var _publishMsgChannel chan kafkaJsonMessage

// _run is a flag that indicates whether the Kafka producer is running.
var _run bool = true

// _bootstrapServers represents the Kafka bootstrap servers.
var _bootstrapServers string = ""

// _enablePublish is a flag that indicates whether the Kafka producer is enabled.
var _enablePublish bool = true

// _producerId represents the ID of the Kafka producer.
var _producerId string = ""

// _consumerGroup represents the consumer group of the Kafka consumer.
var _consumerGroup string = ""

/**
 * Init initializes the kafkaclient instance with the broker string and enables the producer by default.
 * It saves the broker configuration and starts producing messages from the channel.
 * The function creates a Kafka producer and sets up a delivery report handler for produced messages.
 * It also creates a goroutine that continuously waits for messages to be published.
 * The function returns an error if there is an issue initializing the Kafka producer.
 *
 * @param brokers The Kafka broker string.
 * @param enablePublish Flag to enable or disable the Kafka producer.
 * @param producerId The ID of the Kafka producer.
 * @param consumerGroup The consumer group of the Kafka consumer.
 * @return An error if there is an issue initializing the Kafka producer.
 */
func Init(brokers string, enablePublish bool, producerId string, enableDeliveryReport bool, consumerGroup string) error {
	// Todo: should check if already init,
	// Todo: should allow setup some importance config (for producer and consumer), that can be change in time being

	// Save broker config
	_bootstrapServers = brokers
	_enablePublish = enablePublish
	_producerId = producerId
	_consumerGroup = consumerGroup
	if _enablePublish {
		// start produce message from channel
		_publishMsgChannel = make(chan kafkaJsonMessage)
		go func() {
			// Create producer
			kafkaProducer, err := confluentKafka.NewProducer(&confluentKafka.ConfigMap{
				"bootstrap.servers": _bootstrapServers,
				"client.id":         _producerId,
				"acks":              "all",
				//"transactional.id": producerId,
			})
			if err != nil {
				fmt.Printf("[KafkaClient] Error initializing kafka producer client: %v\n", err)
				return
			}
			defer kafkaProducer.Close()
			// Delivery report handler for produced messages
			if enableDeliveryReport {
				go func() {
					for e := range kafkaProducer.Events() {
						// Checking for stop
						if !_run {
							break
						}

						switch ev := e.(type) {
						case *confluentKafka.Message:
							// The message delivery report, indicating success or
							// permanent failure after retries have been exhausted.
							// Application level retries won't help since the client
							// is already configured to do that.
							if ev.TopicPartition.Error != nil {
								fmt.Printf("[KafkaClient-DeliveryReport-%s] Delivery failed: %v\n", ev.TopicPartition, time.Now().Format(time.RFC3339Nano))
							} else {
								fmt.Printf("[KafkaClient-DeliveryReport-%s] Delivered message to %v\n", ev.TopicPartition, time.Now().Format(time.RFC3339Nano))
							}
						case confluentKafka.Error:
							// Generic client instance-level errors, such as
							// broker connection failures, authentication issues, etc.
							//
							// These errors should generally be considered informational
							// as the underlying client will automatically try to
							// recover from any errors encountered, the application
							// does not need to take action on them.
							fmt.Printf("[KafkaClient-DeliveryReport-%s] Error: %v\n", ev, time.Now().Format(time.RFC3339Nano))
						default:
							fmt.Printf("[KafkaClient-DeliveryReport-%s] Ignored event: %s\n", ev, time.Now().Format(time.RFC3339Nano))
						}
					}
				}()
			}
			// Here we create go-routine and continous waiting for message to publish
			// Message is collect from channel "publishMsgChannel"
			msgBatch := make([]kafkaJsonMessage, 0)
			t := time.NewTicker(100 * time.Millisecond)
			defer t.Stop()

			fmt.Println("[KafkaClient] Finish setup producer, start waiting for message to be publish !!!")
			for _run {
				select {
				case <-t.C:
					// Send all message in queue periodically
					if len(msgBatch) > 0 {
						for _, record := range msgBatch {
							err := kafkaProducer.Produce(&confluentKafka.Message{
								TopicPartition: confluentKafka.TopicPartition{Topic: &record.Topic, Partition: confluentKafka.PartitionAny},
								Value:          []byte(record.Message)},
								nil, // delivery channel
							)
							if err != nil {
								if err.(confluentKafka.Error).Code() == confluentKafka.ErrQueueFull {
									// Producer queue is full, wait 1s for messages
									// to be delivered then try again.
									fmt.Println("[KafkaClient] Producer queue is full, wait 1s for messages to be delivered then try again !!!")
									time.Sleep(time.Second)
									continue
								}
								fmt.Printf("[KafkaClient] Failed to produce message: %v\n", err)
							}
						}

						// clear list if message send successfully
						msgBatch = make([]kafkaJsonMessage, 0)
					}
				case msg := <-_publishMsgChannel:
					// add message to queue
					msgBatch = append(msgBatch, msg)
				}
			}
		}()

		return nil
	}

	return nil
}

/**
 * SendJsonMessages publishes a JSON message to the specified topic.
 * It adds the message to the publish message channel for the Kafka producer to consume.
 *
 * @param jsonMessage The JSON message to be sent.
 * @param topic The topic to publish the message to.
 */
func SendJsonMessages(jsonMessage string, topic string) {
	if _publishMsgChannel != nil {
		msg := kafkaJsonMessage{
			Message: jsonMessage,
			Topic:   topic,
		}
		_publishMsgChannel <- msg
	}
}

/**
 * CreateNewConsumer creates a new Kafka consumer and subscribes to the specified topics.
 * It returns the consumer name if successful, otherwise an error is returned.
 * The function checks if the input topic names are valid and not empty.
 * It also checks if the consumer is already subscribed to any of the specified topics.
 * If the consumer name already exists, an error is returned.
 * The function uses the bootstrap servers and consumer group specified in the configuration.
 * It sets additional configuration properties for the consumer, such as session timeout,
 * auto offset reset, and auto offset store.
 * After creating the consumer, it subscribes to the specified topics and adds the consumer
 * to the consumer list.
 *
 * @param topics The topics to subscribe to.
 * @return The consumer name if successful, otherwise an error is returned.
 */
func CreateNewConsumer(topics ...string) (string, error) {
	if len(topics) == 0 {
		return "", errors.New("empty input topic name")
	}
	// Check if all topic name in topics is not valid or empty
	for _, topic := range topics {
		if topic == "" {
			return "", errors.New("invalid input topic name")
		}
	}
	// Check if one of topic is already subscribed before
	for key := range _consumerList {
		for _, topic := range topics {
			if strings.Contains(key, topic) {
				return "", fmt.Errorf("already subscribed to topic %s", topic)
			}
		}
	}
	// Create consumer ID containing all topic names
	consumerName := topicsToConsumerName(topics)
	// Check if this consumerID is already existed
	if _, exists := _consumerList[consumerName]; exists {
		return "", fmt.Errorf("consumer with this name already existed: %s", consumerName)
	}
	// Create consumer and sub to topic
	config := &confluentKafka.ConfigMap{
		"bootstrap.servers": _bootstrapServers,
		"group.id":          _consumerGroup,
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family": "v4",
		"session.timeout.ms":    consumerSessionTimeoutms,
		// Start reading from the first message of each assigned
		// partition if there are no previously committed offsets
		// for this group.
		"auto.offset.reset": "earliest",
		// Whether or not we store offsets automatically.
		"enable.auto.offset.store": true,
		"enable.auto.commit":       true,
	}
	// Create consumer
	kafkaConsumer, err := confluentKafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("[KafkaClient] failed to create new Kafka consumer, error: %v\n", err)
		return "", err
	}
	err = kafkaConsumer.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Printf("[KafkaClient] failed to subscribe to topic %s, error: %v\n", consumerName, err)
		return "", err
	}
	// Add comsumer to list
	_consumerList[consumerName] = kafkaConsumer
	fmt.Printf("[KafkaClient] subscribed to topic %s, start fetching messages!\n", consumerName)
	return consumerName, nil
}

/**
 * ConsumerReadMessage retrieves messages from the given consumer.
 * It reads a message from the Kafka consumer and returns the JSON message, topic, timestamp, and any error encountered.
 * If there are no messages within the timeout, it returns an empty message and no error.
 * If there is an error reading the message, it closes the consumer, removes it from the consumer list,
 * waits for a few seconds, and attempts to reconnect to the topic.
 *
 * @param consumerName The name of the consumer.
 * @return The JSON message, topic, timestamp, and any error encountered.
 */
func ConsumerReadMessage(consumerName string) (jsonMessage string, topic string, timestamp time.Time, err error) {
	if consumerName == "" {
		return "", "", time.Time{}, errors.New("invalid input consumer name")
	}
	kafkaConsumer, exists := _consumerList[consumerName]
	if !exists {
		return "", "", time.Time{}, fmt.Errorf("consumer %s not found", consumerName)
	}
	msg, err := kafkaConsumer.ReadMessage(1000 * time.Millisecond)
	if err != nil {
		if kafkaErr, ok := err.(confluentKafka.Error); ok && kafkaErr.Code() == confluentKafka.ErrTimedOut {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			return "", "", time.Time{}, nil // No message received within the timeout
		}
		fmt.Printf("[KafkaClient] consumer error: %v, name: %s\n", err, consumerName)
		// Trying to subscribe again
		// Close and remove this consumber from map
		kafkaConsumer.Close()
		delete(_consumerList, consumerName)
		// Wait before attempting to reconnect
		time.Sleep(time.Second * 5)
		// Connect to topic again
		topics := consumerNameToTopics(consumerName)
		_, reconnectErr := CreateNewConsumer(topics...)
		if reconnectErr != nil {
			fmt.Printf("[KafkaClient] failed to reconnect consumer %s, error: %v\n", consumerName, reconnectErr)
		}
		return "", "", time.Time{}, err
	}
	return string(msg.Value), *msg.TopicPartition.Topic, msg.Timestamp, nil
}

/**
 * topicsToConsumerName converts a slice of topics into a unique consumer name.
 * It joins the topics using the ";" delimiter and returns the resulting consumer name as a string.
 *
 * @param topics The topics to convert.
 * @return The resulting consumer name.
 */
func topicsToConsumerName(topics []string) string {
	return strings.Join(topics, ";")
}

/**
 * consumerNameToTopics splits the consumerName string into multiple topics.
 * It splits the consumerName using the ";" delimiter and returns the resulting topics as a slice of strings.
 * The original consumerName string remains unchanged.
 *
 * @param consumerName The consumer name to split.
 * @return The resulting topics.
 */
func consumerNameToTopics(consumerName string) []string {
	return strings.Split(consumerName, ";")
}

/**
 * Stop stops the Kafka producer and closes all Kafka consumers.
 */
func Stop() {
	_run = false
	// Should close all consumer
	for topic, consumer := range _consumerList {
		consumer.Close()
		delete(_consumerList, topic)
		fmt.Printf("[KafkaClient] closed consumer for topic %s\n", topic)
	}
}
