# Warning: This Makefile is never used. It's added just for the sake of completeness of this challenge.
# I have populated this with the migrations commands however, this can be easily done in the docker container itself
# by placing the sql files in the /docker-entrypoint-initdb.d/ directory
# Check the Dockerfile under `migrations` directory

SHELL := /bin/bash -eo pipefail
.SILENT:

##@ sqls
.PHONY:sqls
sqls:
	psql -h localhost -U postgres -a -f ./sqls/001-init.sql
	psql -h localhost -U postgres -a -f ./sqls/002-add-time.sql

