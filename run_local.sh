#!/bin/bash

# Функция для остановки всех запущенных процессов
stop_processes() {
  # Останавливаем все запущенные процессы
  docker-compose -f docker-compose-kafka.yml down
  pkill -P $$
}

# Обработчик сигнала SIGINT (Ctrl+C)
trap 'stop_processes; exit 130' INT

source .env
export $(grep -v '^' .env | xargs)

# Запуск Kafka и приложения
docker-compose -f docker-compose-kafka.yml up -d

# Ожидание запуска Kafka и приложения
sleep 5

# Запуск агентов в цикле
for ((i=0; i<$COUNT_AGENTS; i++)); do
  http_port=$(($HTTP_SERVER_PORT+i+1))
  agent_id=$i
  go run -v ./back-end/agent/cmd/app/main.go s --http_port $http_port --agent_id $agent_id &
done

go run -v ./back-end/orkestrator/cmd/app/main.go serve --local --count_agents $COUNT_AGENTS
