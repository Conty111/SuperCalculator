#!/bin/bash

# Чтение количества агентов из JSON файла
agents_count=$(jq '.agents | length' agents.json)

# Замена значения переменной AGENTS_COUNT в docker-compose файле
export AGENTS_COUNT=$agents_count
docker build -t svc-agent:latest DockerfileAgent .
# Запуск docker-compose (все компоненты, кроме agent)
docker-compose up orkestrator -d

for ((i=0;i<agents_count;i++))
do
    docker run -d \
            --name agent$i \
            --restart on-failure \
            --env AGENTS_COUNT=$agents_count \
            --env-file enviroments/agent.env \
            --env-file enviroments/kafka.env \
            --env-file enviroments/.env \
            --network svc-network \
            svc-agent \
            /app serve $i
done