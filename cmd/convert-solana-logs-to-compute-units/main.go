package main

import (
	"encoding/json"
	"io"
	"log"

	"github.com/fluidity-money/devops-hiring-challenge/lib"
)

const (
	// TopicSolanaSlots to read from to get the Solana slots
	TopicSolanaSlots = devops_hiring_challenge.TopicSolanaSlots

	// TopicAverageComputeUnits to send the average compute units down
	TopicAverageComputeUnits = devops_hiring_challenge.TopicAverageComputeUnits
)

func main() {
	err := devops_hiring_challenge.GetMessages(TopicSolanaSlots, func(reader io.Reader) {
		var solanaWebsocketLog devops_hiring_challenge.SolanaWebsocketLog

		if err := json.NewDecoder(reader).Decode(&solanaWebsocketLog); err != nil {
			log.Fatalf(
				"Failed to decode a message off Kafka for the websocket log! %v",
				err,
			)
		}

		log.Printf(
			"Got this message off Kafka: %v",
			solanaWebsocketLog,
		)

		computes := extractComputes(solanaWebsocketLog)

		summedUnits := computeSummedUnits(computes)

		log.Printf(
			"Summed compute units for this message is %v!",
			summedUnits,
		)

		err := devops_hiring_challenge.PublishMessage(
			TopicAverageComputeUnits,
			summedUnits,
		)

		if err != nil {
			log.Fatalf(
				"Failed to send down the Kafka the summed computes from the logs here! %v",
				err,
			)
		}
	})

	log.Fatalf(
		"Solana slots consumption from Kafka topic %#s ended preemptively! %v",
		TopicSolanaSlots,
		err,
	)
}
