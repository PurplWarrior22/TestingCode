#!/bin/bash
#generates the server boilerplate 

oapi-codegen -package swagger ../../api-docs/madden-swagger.yaml > swagger/madden.gen.go