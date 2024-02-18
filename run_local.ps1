# Функция для остановки всех запущенных процессов
function Stop-Processes {
    # Останавливаем все запущенные процессы
    docker-compose -f docker-compose-kafka.yml down
    Stop-Process -Id $PID
}

# Обработчик сигнала SIGINT (Ctrl+C)
$PID = $PID

trap {
    Stop-Processes
    exit 130
} -signal INT

# Импорт переменных окружения из файла .env
$envVars = Get-Content ".env" | ForEach-Object { $_ -split '=' }
foreach ($envVar in $envVars) {
    $envVarName = $envVar[0].Trim()
    $envVarValue = $envVar[1].Trim()
    [System.Environment]::SetEnvironmentVariable($envVarName, $envVarValue, [System.EnvironmentVariableTarget]::Process)
}

# Запуск Kafka и приложения
docker-compose -f docker-compose-kafka.yml up -d

# Ожидание запуска Kafka и приложения
Start-Sleep -Seconds 5

# Запуск агентов в цикле
for ($i=0; $i -lt $env:COUNT_AGENTS; $i++) {
    $http_port = $env:HTTP_SERVER_PORT + $i + 1
    $agent_id = $i
    Start-Process "go" -ArgumentList "-v", "./back-end/agent/cmd/app/main.go", "s", "--http_port", $http_port, "--agent_id", $agent_id -NoNewWindow
}

Start-Process "go" -ArgumentList "-v", "./back-end/orkestrator/cmd/app/main.go", "serve", "--local", "--count_agents", $env:COUNT_AGENTS -NoNewWindow

# Ожидание завершения всех процессов
Wait-Process -Name "go"
