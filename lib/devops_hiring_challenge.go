package devops_hiring_challenge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/segmentio/kafka-go"
)

// EnvKafkaLeaderAddress to point to connect to the leader for this cluster
const EnvKafkaLeaderAddress = `FLU_KAFKA_LEADER`

const (
	// TopicSolanaSlots contains the received logs from Solana
	TopicSolanaSlots = `hiring.solana.slots`

	// TopicAverageComputeUnits contains the average compute units from the
	// Solana logs received
	TopicAverageComputeUnits = `hiring.problem.average-compute-units`
)

type (
	requestKafkaWriter struct {
		topic    string
		response chan *kafka.Writer
	}

	requestKafkaReader struct {
		topic    string
		response chan *kafka.Reader
	}
)

var (
	chanKafkaWriters  = make(chan requestKafkaWriter, 0)
	chanKafkaReaders  = make(chan requestKafkaReader, 0)
	chanKafkaShutdown = make(chan bool, 0)
)

// PublishMessage by using a lookup to get the writer each time if the
// topic isn't found already, then sending down a writer to Kafka
func PublishMessage(topicName string, message any) (err error) {
	chanWriter := make(chan *kafka.Writer, 0)

	chanKafkaWriters <- requestKafkaWriter{
		topic:    topicName,
		response: chanWriter,
	}

	writer := <-chanWriter

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(message); err != nil {
		return fmt.Errorf(
			"failed to write JSON to a buffer before sending! %v",
			err,
		)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: buf.Bytes(),
	})

	if err != nil {
		return fmt.Errorf(
			"failed to write messages to Kafka! %v",
			err,
		)
	}

	return
}

// GetMessages from Kafka using the topic name and calling the function
// each time
func GetMessages(topicName string, f func(message io.Reader)) error {
	chanReader := make(chan *kafka.Reader, 0)

	chanKafkaReaders <- requestKafkaReader{
		topic:    topicName,
		response: chanReader,
	}

	reader := <-chanReader

	for {
		message, err := reader.ReadMessage(context.Background())

		if err != nil {
			return fmt.Errorf(
				"failed to read a message off Kafka topic %#v! %v",
				topicName,
				err,
			)
		}

		var buf bytes.Buffer

		_, _ = buf.Write(message.Value)

		f(&buf)
	}
}

func Close() {
	chanKafkaShutdown <- true
}

func init() {
	kafkaLeaderAddress := os.Getenv(EnvKafkaLeaderAddress)

	if kafkaLeaderAddress == "" {
		log.Fatalf("%s is not set!", EnvKafkaLeaderAddress)
	}

	kafkaLeaderUrl, err := net.ResolveTCPAddr("tcp", kafkaLeaderAddress)

	if err != nil {
		log.Fatalf(
			"Failed to parse %#v as a TCP address! %v",
			kafkaLeaderUrl,
			err,
		)
	}

	go func() {
		var (
			writerPool = make(map[string]*kafka.Writer, 2)
			readerPool = make(map[string]*kafka.Reader, 2)
		)

		defer func() {
			for topic, writer := range writerPool {
				if err := writer.Close(); err != nil {
					log.Printf(
						"A writer at topic %#s that was closing errored out %v!",
						topic,
						err,
					)
				}
			}

			for topic, reader := range readerPool {
				if err := reader.Close(); err != nil {
					log.Printf(
						"A reader at topic %#s that was closing errored out %v!",
						topic,
						err,
					)
				}
			}
		}()

		for {
			select {
			case request := <-chanKafkaWriters:
				var (
					topic    = request.topic
					response = request.response
				)

				_, ok := writerPool[topic]

				if !ok {
					log.Printf(
						"Creating a new writer for the topic %#v!",
						topic,
					)

					writerPool[topic] = &kafka.Writer{
						Addr:  kafkaLeaderUrl,
						Topic: topic,
						Async: true,
					}
				}

				response <- writerPool[topic]

			case request := <-chanKafkaReaders:
				var (
					topic    = request.topic
					response = request.response
				)

				_, ok := readerPool[topic]

				if !ok {
					log.Printf(
						"Creating a new reader for the topic %#v!",
						topic,
					)

					readerPool[topic] = kafka.NewReader(kafka.ReaderConfig{
						Brokers: []string{kafkaLeaderAddress},
						Topic:   topic,
					})
				}

				response <- readerPool[topic]

			case shutdownRequest := <-chanKafkaShutdown:
				if !shutdownRequest {
					continue
				}

				break
			}
		}
	}()
}
