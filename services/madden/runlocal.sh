#!/bin/bash
#runs against the default local website build docker-compose-with-elasticsearch

export DB_USERNAME=developer
export DB_PASSWORD=development
export DB_HOST=localhost
export DB_PORT=9876
export IMAGE_PATH=http://www.google.com/

go run .