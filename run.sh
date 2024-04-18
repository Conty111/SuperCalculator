#!/bin/bash

# Чтение количества агентов из JSON файла
agents_count=$(jq '.agents | length' system_config_docker.json)

# Замена значения переменной AGENTS_COUNT в docker-compose файле
export AGENTS_COUNT=$agents_count
#echo $agents_count
if [ -z "$(docker images -q svc-agent:script 2> /dev/null)" ]; then
  docker build -t svc-agent:script -f DockerfileAgent .
fi
# Запуск docker-compose (все компоненты, кроме agent)
docker compose up -d orkestrator

for ((i=1;i<agents_count+1;i++))
do
    grpc_port=$(jq -r --arg index "$((i-1))" '.agents[$index | tonumber] | .grpc_port' system_config_docker.json)
    http_port=$(jq -r --arg index "$((i-1))" '.agents[$index | tonumber] | .http_port' system_config_docker.json)
    docker run -d \
            --name agent$i \
            --add-host=agent$i:0.0.0.0 \
            --restart on-failure \
            -p "$grpc_port:$grpc_port" \
            -p "$http_port:$http_port" \
            --env-file=enviroments/agent.env \
            --env-file=enviroments/kafka.env \
            --env-file=enviroments/.env \
            --network=supercalculator_calculator-network \
            svc-agent:script \
            /app serve $((i-1))
#            yes
done

#docker-compose up -d orkestrator