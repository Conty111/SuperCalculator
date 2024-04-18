# Остановка и удаление всех запущенных контейнеров
docker ps -aq --filter "name=^/agent" | ForEach-Object { docker stop $_ }
docker ps -aq --filter "name=^/agent" | ForEach-Object { docker rm $_ }
docker-compose down
