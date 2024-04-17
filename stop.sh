#!/bin/bash

# Остановка и удаление всех запущенных контейнеров
docker stop $(docker ps -aq --filter "name=^/agent")
docker rm $(docker ps -aq --filter "name=^/agent")
docker-compose down
docker network rm svc-network
