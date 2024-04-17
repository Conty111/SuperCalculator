#!/bin/bash

# Остановка и удаление всех запущенных контейнеров
docker stop $(docker ps -aq --filter "name=^/agent")
docker rm $(docker ps -aq --filter "name=^/agent")
docker rmi $(docker images -q --filter "name=svc-agent")
docker-compose down
