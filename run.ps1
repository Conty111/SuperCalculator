# Чтение количества агентов из JSON файла
$agents_count = (Get-Content -Raw -Path "system_config_docker.json" | ConvertFrom-Json).agents.Count

# Замена значения переменной AGENTS_COUNT в docker-compose файле
$env:AGENTS_COUNT = $agents_count

# Проверка наличия образа в Docker
if (-not (docker images -q svc-agent:script 2> $null)) {
    docker build -t svc-agent:script -f DockerfileAgent .
}

# Запуск docker-compose (все компоненты, кроме agent)
docker compose up -d orkestrator

for ($i=1; $i -le $agents_count; $i++) {
    $agent = (Get-Content -Raw -Path "system_config_docker.json" | ConvertFrom-Json).agents[$i - 1]
    $grpc_port = $agent.grpc_port
    $http_port = $agent.http_port

    docker run -d `
        --name agent$i `
        --add-host=agent$i:0.0.0.0 `
        --restart on-failure `
        -p "$grpc_port:$grpc_port" `
        -p "$http_port:$http_port" `
        --env-file=enviroments/agent.env `
        --env-file=enviroments/kafka.env `
        --env-file=enviroments/.env `
        --network=supercalculator_calculator-network `
        svc-agent:script `
        /app serve $i
}