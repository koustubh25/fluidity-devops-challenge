package main

import (
	"log"
	"os"

	"github.com/fluidity-money/devops-hiring-challenge/lib"

	"github.com/gorilla/websocket"
)

// TopicSolanaSlots to send down with messages from Solana
const TopicSolanaSlots = devops_hiring_challenge.TopicSolanaSlots

// EnvSolanaWebsocketUrl to use to connect to Solana upstream
const EnvSolanaWebsocketUrl = `FLU_SOLANA_WEBSOCKET_URL`

func main() {
	solanaWebsocketUrl := os.Getenv(EnvSolanaWebsocketUrl)

	if solanaWebsocketUrl == "" {
		log.Fatalf("%s is not set!", EnvSolanaWebsocketUrl)
	}

	defer devops_hiring_challenge.Close()

	log.Printf(
		"Connecting to the Solana websocket at %#s...",
		solanaWebsocketUrl,
	)

	client, _, err := websocket.DefaultDialer.Dial(solanaWebsocketUrl, nil)

	if err != nil {
		log.Fatalf(
			"Failed to connect to %s! %v",
			EnvSolanaWebsocketUrl,
			err,
		)
	}

	log.Printf(
		"Connected to the Solana websocket at %#s",
		solanaWebsocketUrl,
	)

	defer client.Close()

	writerHello, err := client.NextWriter(websocket.TextMessage)

	if err != nil {
		log.Fatalf(
			"Failed to open a writer to write to the Solana websocket the subscribe message! %v",
			err,
		)
	}

	if err := websocketWriteSubscribeMessage(writerHello); err != nil {
		log.Fatalf(
			"Failed to write the subscribe message to the websocket! %v",
			err,
		)
	}

	for {
		var solanaWebsocketLog devops_hiring_challenge.SolanaWebsocketLog

		log.Printf(
			"Waiting for a message on the Solana websocket...",
		)

		if err := client.ReadJSON(&solanaWebsocketLog); err != nil {
			log.Fatalf(
				"Failed to read a JSON message off the Solana websocket! %v",
				err,
			)
		}

		log.Printf(
			"Received this message from the Solana websocket deserialised as %v",
			solanaWebsocketLog,
		)

		err := devops_hiring_challenge.PublishMessage(
			TopicSolanaSlots,
			solanaWebsocketLog,
		)

		if err != nil {
			log.Fatalf(
				"Failed to send a message down the Kafka hiring queue! %v",
				err,
			)
		}
	}
}
