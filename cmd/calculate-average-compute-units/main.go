package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fluidity-money/devops-hiring-challenge/lib"

	_ "github.com/lib/pq"
)

// TopicComputeUnits to get the messages off the kafka from
const TopicComputeUnits = devops_hiring_challenge.TopicAverageComputeUnits

// EnvTimescaleUri to connect to the proper Timescale database!
const EnvTimescaleUri = `FLU_TIMESCALE_URI`

func main() {
	timescaleUri := os.Getenv(EnvTimescaleUri)

	if timescaleUri == "" {
		log.Fatalf("%s is not set!", EnvTimescaleUri)
	}

	defer devops_hiring_challenge.Close()

	database, err := sql.Open("postgres", timescaleUri)

	if err != nil {
		log.Fatalf(
			"Failed to open the connection to the Timescale database! %v",
			err,
		)
	}

	log.Printf("Connected to the Timescale database!")

	defer database.Close()

	err = devops_hiring_challenge.GetMessages(TopicComputeUnits, func(message io.Reader) {
		var computeUnits devops_hiring_challenge.ComputeUnits

		if err := json.NewDecoder(message).Decode(&computeUnits); err != nil {
			log.Fatalf(
				"Failed to decode the compute units in a message off Kafka! %v",
				err,
			)
		}

		log.Printf(
			"Received a message with the compute units as %v",
			computeUnits,
		)

		_, err = database.Exec(
			"INSERT INTO average_compute_units (compute_units) VALUES ($1);",
			computeUnits,
		)

		if err != nil {
			log.Fatalf(
				"Failed to insert into the timescale database the compute units! %v",
				err,
			)
		}

		log.Printf("Inserted into the database the compute units!")

		row := database.QueryRow(
			`SELECT AVG(compute_units) FROM average_compute_units`,
		)

		if err := row.Err(); err != nil {
			log.Fatalf(
				"Failed to get the average compute units from the database! %v",
				err,
			)
		}

		var averageComputeUnits float64

		if err := row.Scan(&averageComputeUnits); err != nil {
			log.Fatalf(
				"Failed to scan the compute units from the database! %v",
				err,
			)
		}

		fmt.Println("Average compute units: %v", averageComputeUnits)
	})
}
