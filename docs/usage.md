# SuperCalculator

## Table of contents

- ### [What is it?](..%2FREADME.md#what-is-it)
- ### [How to run?](..%2FREADME.md#running)
- ### [How to use?](usage.md)
- ### [How it works?](howItWorks.md)
- ### [Project sructure](..%2FREADME.md#project-structure)

## How to use
* Сразу к curl запросам: **[Requests](#curl-requests)**

По умолчанию, при перезапуске системы в Docker Compose данные НЕ сохраняются.\
Если хотите, чтобы данные сохранялись, раскомментируйте строки volume для db в [docker-compose.yml](..%2Fdocker-compose.yml )
```
  db:
    image: postgres:latest
    container_name: db
    env_file:
      - enviroments/docker.db.env
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=postgres
    ports:
      - "5433:5432"
#    volumes:                                                   <- эту строку
#      - ./back-end/db/data/postgres:/var/lib/postgresql/data   <- и эту строку
    networks:
      - calculator-network
```
Основные настройки системы находятся в env файлах в [enviroments](..%2Fenviroments) директории. \
Также настройки системы задаются в [system_config.json](..%2Fsystem_config.json) 
(и в [system_config_docker.json](..%2Fsystem_config_docker.json) для Docker Compose).
#### Структура system_config.json
```
{
  "brokers": [
    {"address": "localhost", "port":  29092} - адреса брокеров для подключения
  ],
  "agents": [ - массив всех агентов 
   {
      "name": "agent1", - просто имя агента
      "address": "0.0.0.0",
      "broker_partition": 0, - номер партиции в брокере (уникальное значение)
      "consumer_group": "agents", - название группы consumer-ов в брокере
      "broker_commit_interval": 1,
      "http_port": 8001, - порт прослушивания HTTP запросов
      "grpc_port": 8002 - порт прослушивания GRPC запросов
    },
    {
      "name": "agent2",
      "address": "0.0.0.0",
      "broker_partition": 1,
      "consumer_group": "agents",
      "broker_commit_interval": 1,
      "http_port": 8003,
      "grpc_port": 8004
    }
  ]
}
```
Подразумевается, что все указанные агенты в system_config.json должны использоваться.\
Никакие из портов не должны повторяться, добавить нового агента можно просто скопировав предыдущего и указав новые уникальные значения

## Users
В системе присутствует хранение пользователей и их разграничение. Также присутствует разграничение прав доступа:

Есть роли **admin** и **common**. \
Admin может выполнять все запросы, он не высвечивается при выводе всех пользователей.\
Common не может удалять или редактировать пользователей.

Создать admin пользователя можно только вручную в БД, поэтому для удобного создания admin пользователя можно воспользоваться
```
make admin-postgres-docker
```
или
```
make admin-sqlite-local
```
И тогда создастся пользователь admin с данными:\
* email: admin@mail.ru
* password: 12345

## Tests
Для запуска unit-тестов можно использовать ```make test-unit``` или ```go test -v -cover ./...```\
Модульные тесты были написаны в основном для 2 модулей:
- Agent - вычисление математических выражений (23 теста, [calculator_test.go](..%2Fback-end%2Fagent%2Finternal%2Fservices%2Fcalculator_test.go))
- Orkestrator - контроллер для работы с пользователями (18 тестов, [controller_test.go](..%2Fback-end%2Forkestrator%2Finternal%2Ftransport%2Fweb%2Fcontrollers%2Fapiv1%2Fuser%2Fcontroller_test.go))

Остальную логику необходимо тестировать в рамках интеграционных тестов.\
На данный интеграционные тесты отсутствуют

## Requests

Если вы пользуетесь **Postman**, то можно импортировать коллекцию 
[Calculator.postman_collection.json](Calculator.postman_collection.json) и использовать её.

### Curl requests
Register
```
curl -L -m 5 'http://localhost:8000/api/v1/users/register' \
-H 'Content-Type: application/json' \
--data-raw '{
    "username": "testUser",
    "email": "testmail@mail.ru",
    "password": "123"
}'
```
Login
```
curl -L -m 5 'http://localhost:8000/api/v1/users/login' \
-H 'Content-Type: application/json' \
--data-raw '{
    "email": "admin@mail.ru",
    "password": "12345"
}
'
```
Create task
```
curl -L -m 5 'http://localhost:8000/api/v1/tasks' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer <JWT токен, полученный из Login или Register>' \
-d '{
    "expression": "111+11+22"
}'
```
Get all tasks
```
curl -L -m 5 'http://localhost:8000/api/v1/tasks' \
-H 'Authorization: Bearer <JWT токен, полученный из Login или Register>'
```
Set settings
* "time_retry" в секундах
* "timeout_response" - в секундах
* "add_time" - в миллисекундах
* "division_time" - в миллисекундах
* "subtract_time" - в миллисекундах
* "multiply_time" - в миллисекундах
```
curl -L -m 5 -X PUT 'http://localhost:8000/api/v1/workers/settings' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer <JWT токен, полученный из Login или Register>'
-d '{
    "time_retry": 30, 
    "timeout_response": 5,
    "add_time": 5000,
    "division_time": 500,
    "subtract_time": 500,
    "multiply_time": 500
}'
```
Get workers info
```
curl -L -m 5 'http://localhost:8000/api/v1/workers/info' \
-H 'Authorization: Bearer <JWT токен, полученный из Login или Register>'
```
Get my info
```
curl -L -m 5 'http://localhost:8000/api/v1/users/me' \
-H 'Authorization: Bearer <JWT токен, полученный из Login или Register>'
```
Get all users
```
curl -L -m 5 'http://localhost:8000/api/v1/users' \
-H 'Authorization: Bearer <JWT токен, полученный из Login или Register>'
```
Delete user (в URL указывается ID пользователя)
```
curl -L -m 5 -X DELETE 'http://localhost:8000/api/v1/users/1' \
-H 'Authorization: Bearer <JWT токен, полученный из Login или Register>'
```
Update user (в URL указывается ID пользователя)
```
curl -L -m 5 -X PATCH 'http://localhost:8000/api/v1/users/1' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer <JWT токен, полученный из Login или Register>'
-d '{
    "param": "Username",
    "value": "NewTestUser"
}'
```
Get orkestrator status
```
curl -L -m 5 'http://localhost:8000/api/v1/status'
```