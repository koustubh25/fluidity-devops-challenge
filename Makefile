SHELL := /bin/bash -eo pipefail
.SILENT:

POSTGRES_PASSWORD ?= postgres # default password for timescaledb

# start kafka server along with zookeeper
.PHONY:kafka
kafka: stop
	docker-compose run kafka

# start timescaledb with the initial migrations
.PHONY:timescaledb
timescaledb: stop
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) docker-compose run timescaledb

# start solana logs connector, this will also start kafka along with zookeeper
.PHONY:connector
connector: stop
	docker-compose run solana-connector

# start solana logs connector, this will also start kafka along with zookeeper
.PHONY:calculate
calculate: stop
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) docker-compose run calculate-avg-compute-units

.PHONY:challenge
challenge: stop
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) docker-compose up --build


# stop all
.PHONY:stop
stop:
	docker-compose down

# clean everything
.PHONY:clean
clean: stop
	docker-compose down
	docker volume rm -f $(shell docker volume ls -q  | grep -i timescale-db)
	docker rmi -f $(shell docker images --filter=reference='*convert*:*' -q)
	docker rmi -f $(shell docker images --filter=reference='*solana*:*' -q)
	docker rmi -f $(shell docker images --filter=reference='*timescaledb*:*' -q)
	docker rmi -f $(shell docker images --filter=reference='*calculate*:*' -q)


