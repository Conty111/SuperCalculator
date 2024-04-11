#!/bin/bash

# Функция для остановки всех запущенных процессов
stop_processes() {
  # Останавливаем все запущенные процессы
  docker compose -f docker-compose-kafka.yml down
  pkill -P $$
}

export COUNT_AGENTS=$(jq '.agents | length' config.json)


# Обработчик сигнала SIGINT (Ctrl+C)
trap 'stop_processes; exit 130' INT

cp .env.example .env
source .env
source kafka.env
source agent.env
source orkestrator.env
export $(grep '^' .env | xargs)
export $(grep '^' orkestrator.env | xargs)
export $(grep '^' agent.env | xargs)
export $(grep '^' kafka.env | xargs)

# Запуск Kafka и приложения
docker-compose -f docker-compose-kafka.yml up -d

# Ожидание запуска Kafka и приложения
sleep 3

# Запуск агентов в цикле
for ((i=0; i<$COUNT_AGENTS; i++)); do
  go run -v ./back-end/agent/cmd/app/main.go s $i &
done

go run -v ./back-end/orkestrator/cmd/app/main.go s
