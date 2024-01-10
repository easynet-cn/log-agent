package object

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/easynet-cn/log-agent/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var producers map[string]sarama.AsyncProducer

func InitProducer(viper *viper.Viper) {
	producers := make(map[string]sarama.AsyncProducer)

	for k := range viper.GetStringMap("kafka") {
		config := sarama.NewConfig()

		config.Producer.Compression = sarama.CompressionGZIP

		if client, err := sarama.NewClient(viper.GetStringSlice(fmt.Sprintf("kafka.%s.brokers", k)), config); err != nil {
			log.Logger.Panic("kafka", zap.Error(err))

			panic(err)
		} else if producer, err := sarama.NewAsyncProducerFromClient(client); err != nil {
			log.Logger.Panic("kafka", zap.Error(err))

			panic(err)
		} else {
			producers[k] = producer
		}
	}
}

func Send(producerName string, topic string, data []byte) {
	bytes := sarama.ByteEncoder(data)

	if producer, ok := producers[producerName]; ok {
		producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: bytes}
	}
}

func Close() {
	for _, producer := range producers {
		producer.Close()
	}
}
