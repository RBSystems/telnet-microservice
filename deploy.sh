#!/usr/bin/env bash
# Pasted into Jenkins to build (will eventually be fleshed out to work with a Docker Hub and Amazon AWS)

echo "Stopping running application"
docker stop telnet-microservice
docker rm telnet-microservice

echo "Building container"
docker build -t byuoitav/telnet-microservice .

echo "Starting the new version"
docker run -d --restart=always --name telnet-microservice -p 8001:8001 byuoitav/telnet-microservice:latest
