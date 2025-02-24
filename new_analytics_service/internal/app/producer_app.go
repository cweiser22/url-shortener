package app

import "github.com/IBM/sarama"

type ProducerApp struct {
	KafkaProducer sarama.SyncProducer
}
