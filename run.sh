#!/bin/bash

# Чтение количества агентов из JSON файла
agents_count=$(jq '.agents | length' system_config_docker.json)

# Замена значения переменной AGENTS_COUNT в docker-compose файле
export AGENTS_COUNT=3
docker build -t svc-agent:latest -f DockerfileAgent .
# Запуск docker-compose (все компоненты, кроме agent)
docker-compose up -d kafka

for ((i=1;i<agents_count+1;i++))
do
    docker run -d \
            --name agent$i \
            --add-host=agent$i:127.0.0.1 \
            --restart on-failure \
            --env-file=enviroments/agent.env \
            --env-file=enviroments/kafka.env \
            --env-file=enviroments/.env \
            --network=supercalculator_calculator-network \
            svc-agent \
            /app serve ${i-1}
#            yes
done

docker-compose up -d orkestrator