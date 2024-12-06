package mq_demo

import (
	"fmt"
	"github.com/Shopify/sarama"
)

/*	安装：
	go get github.com/Shopify/sarama
*/

func SendMessage() {
	brokers := []string{"127.0.0.1:9092"}
	topic := "testMQ"

	config := sarama.NewConfig()

	producer, err := sarama.NewSyncProducer(brokers, topic)
	defer producer.Close()

	message := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("hello mq"),
	}

	partition, offset, err := producer.SendMessage(message)

}

func consumeMessage() {
	brokers := []string{"127.0.0.1:9092"}
	topic := "testMQ"

	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer(brokers, topic)
	defer consumer.Close()

	partitions, err1 := consumer.Partitions(topic)

	for _, partition := range partitions {
		go func(partition int32) {
			partitionConsumer, err := client.ConsumerPartition(topic, partition, sarama.OffsetNewest)
			defer partitionConsumer.Close()

			for msg := range partitionConsumer.Messages() {
				fmt.Println("consumer message from topic:%s,partition:%d, message is %s", topic, msg.Partition, msg.Value)
			}
		}(partition)
	}
}
