#!/bin/bash

# Прочитать JSON и найти количество агентов
num_agents=$(jq '.agents | length' system_config.json)

# Установить количество агентов в переменную среды
export AGENTS_COUNT=$num_agents
docker compose -f docker-compose.yml up -d
