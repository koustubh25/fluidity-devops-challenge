
# Devops Hiring Challenge

In this ridiculous devops hiring challenge, you are expected to (in no
particular order):

1. Prepare a Terraform configuration for the services

2. Complete partially filled Dockerfiles for the cmdlets based off a
   root container

3. Prepare a single docker-compose.yml file that can auto run
   everything and add a Makescript step that will run everything with
   default environment variables given below

4. Prepare Docker containers for each of the microservices contained
   within `cmd`

5. Load the simple database migration files during the setup for Timescale

## Topology

`solana-logs-connector` receives Solana logs off the websocket
and sends them down Kafka. `convert-solana-logs-to-compute-units`
captures the distilled form of the Solana logs and sends them down
Kafka. `calculate-average-compute-units` uses Timescale, receives the
message of the compute units and calculates the average at that time.

## Infrastructure

Kafka and Timescale will need to be accessible, and the latter will need
to have migrations loaded from the `migrations` directory.

## Environment variables

|             Name           |                                  Description                                   |
|----------------------------|--------------------------------------------------------------------------------|
| `FLU_KAFKA_LEADER`         | The broker to connect to using Kafka.                                          |
| `FLU_SOLANA_WEBSOCKET_URL` | The Solana websocket to connect to (use `https://api.mainnet-beta.solana.com`) |
| `FLU_TIMESCALE_URI`        | URI to use to connect to Timecsale with.                                       |

Good luck!

Reach out to (mailto:alex@fluidity.money)[alex.fluidity.money] if you
need anything.
