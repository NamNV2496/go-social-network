package producer

import (
	"context"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/namnv2496/post-service/configs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client interface {
	Produce(ctx context.Context, queueName string, payload []byte) error
}

type client struct {
	saramaSyncProducer sarama.SyncProducer
}

func newSaramaConfig(clientId string) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Retry.Max = 1
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = clientId
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

func NewClient(
	config configs.Config,
) Client {
	mqConfig := config.Kafka
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = mqConfig.Addresses
	} else {
		log.Println("KAFKA_BROKER environment variable is set: ", kafkaBroker)
	}
	saramaSyncProducer, err := sarama.NewSyncProducer(
		[]string{kafkaBroker},
		newSaramaConfig(mqConfig.ClientID),
	)
	if err != nil {
		log.Panicln("failed to create sarama sync producer: ", err)
		return nil
	}

	return &client{
		saramaSyncProducer: saramaSyncProducer,
	}
}

func (c client) Produce(ctx context.Context, queueName string, payload []byte) error {
	if _, _, err := c.saramaSyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: queueName,
		Value: sarama.ByteEncoder(payload),
	}); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
