#!/bin/bash

# Парсинг флага количества агентов (по умолчанию 2)
while [[ "$#" -gt 0 ]]; do
  case $1 in
    -a|--agent_count)
      num_agents="$2"
      shift 2
      ;;
    *)
      echo "Unknown parameter passed: $1" >&2
      exit 1
      ;;
  esac
done

# Проверка наличия значения флага, и установка значения по умолчанию
if [ -z "$num_agents" ]; then
  num_agents=2
fi

# Запуск Kafka и приложения
docker-compose -f docker-compose-kafka.yml up -d

# Ожидание запуска Kafka и приложения
sleep 5

# Запуск агентов в цикле
for ((i=0; i<$num_agents; i++)); do
  http_port=$((8000+i))
  agent_id=$i
  go run -v ./cmd/app/main.go s --http_port $http_port --agent_id $agent_id &
done

go run -v .back-end/orkestrator/cmd/app/main.go serve --local --base_http_port=8000
