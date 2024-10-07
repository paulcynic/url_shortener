#!/bin/bash

export CONFIG_PATH="../config/local.yaml"
export ENV="local"
# export HTTP_SERVER_PASSWORD="changeit"

go run ./url-shortener
