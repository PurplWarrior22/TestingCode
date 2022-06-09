#!/bin/bash
#script which sets env and runs tests against a running database. assumes database is running on host/ports as specified in the environment variables 
export DB_USERNAME="postgres"
export DB_PASSWORD="development"
export DB_PORT="5432"
export DB_HOST="localhost"
export DB_NAME="postgres"
go test -count=1 ./...