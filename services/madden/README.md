# Madden 

This service handles interactions between HTTP clients and the madden database.

## Dependencies
Project dependencies are available in the go.mod file.

## Configuration 

All configuration as found in [maddendb](../maddendb/README.md) additional configuration below:

### IMAGE_PATH
the path which will be prepended to images before being returned to a client. As an example, if the value "image.png" is stored in the database and this 
is set to http://imageserver.images.com/ the returned path will be "http://imageserver.images.com/image.png"

FORMAT: string

DEFAULT: none 

EXAMPLE: http://imageserver.com

### SERVER_PORT
an optional port the server will listen on 

FORMAT: string

DEFAULT: 8080

EXAMPLE: 8080

## Building
This service is designed to be packaged as a docker image.

A sample build command is:

```
docker build --build-arg USER=$OTH_USER --build-arg PASS=$OTH_PASS -t docker-hfrd.di2e.net/maddenimage:0.1 .

```

### Build arguments description 

Two build arguments are required to authenticate with a private repository.

USER: the username to authenticate to version control with 

PASS: the password or token to authenticate to version control with

## Running The Server
The server can be run using the above built docker images, and supplying appropriate environment variables. 

A script 'runlocal.sh' is supplied, it supplies all required environment configuration to run with the developer docker compose setup.

## oapi-codegen 

This project uses the oapi-codegen swagger generator to build all server boilerplate. A build script (generateserver.sh) is supplied that will update the server based on whatever is found in the api-docs/madden-swagger.yaml file.

Full documentation on this generator can be found [here](https://github.com/deepmap/oapi-codegen)
