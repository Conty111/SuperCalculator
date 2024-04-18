# SuperCalculator

## Table of contents

- ### [Description](#description)
- ### [How to run?](#how-to-run)
- ### [How to use?](docs%2Fusage.md#how-to-use)
- ### [How it works?](docs%2FhowItWorks.md)
- ### [Project sructure](#structure)

## Description

**SuperCalculator** - это распределенный калькулятор или проект, имитирующий отказоустойчивую
систему орекстрации трафика/задач/сообщений между сущностями. 

Также реализован пользовательский функционал (регистрация и авторизация, разграничение доступа на основе ролей, просмотр своих задач, администрирование пользователей и т.п.)

_P. S. front-end часть в проекте на данный момент не работает._

## How to run

### Requirements
Для запуска системы нужен **Go v1.21.6** (и выше) и **Docker** вместе с **Docker Compose**

### Windows preinstalling

Если планируете запускать локально, то Go-sqlite3 требует gcc компилятор, которого по умолчанию нет в Windows. Для установки, можно перейти на сайт https://jmeubank.github.io/tdm-gcc/ и установить его
![image](https://github.com/Conty111/SuperCalculator/assets/90860829/5fed60e6-442f-4ec7-aafb-5360ba3e3e50)

### Running
Запустить можно локально, с помощью Docker Compose или комбинированно (kafka запускается в любом случае в docker compose).

1. С помощью Docker Compose: запустится система с агентами из [system_config_docker.json](system_config_docker.json)

   На **Linux**
    ```
    ./run.sh
    ./stop.sh
    ```
    На **Windows** (с помощью PowerShell)
    ```
    .\run.ps1
    .\stop.ps1
    ```
## Quick start

Для качественного просмотра работы системы, рекомендую пропустить этот пункт и спуститься ниже к пункту **Start in multiple termunals**. Однако, для простого запуска, хватит и этого

1. Установить зависимости``go mod tidy``

#### On Mac/Linux

```
./run_local.sh
```

#### On Windows

```
.\run_local.ps1
```

## Start in multiple terminals

1. Установить зависимости
```
go mod tidy
```
2. Запустить kafka брокер
```
docker-compose -f docker-compose-kafka.yml up -d
```
3. Переименовать .env.example в .env и отредактировать HTTP порты и адреса под себя (указать URL адреса для агентов, например)
```
HTTP_AGENT_ADDRESSES="localhost:8001/api/v1;localhost:8002/api/v1"
```
4. Запустить оркестратор
```
go run -v ./back-end/orkestrator/cmd/app/main.go serve
```
5. Запустить агенты (в отдельных терминалах)
```
go run -v ./back-end/agent/cmd/app/main.go s --http_port <порт агента> --agent_id <id от 0 до макс. кол-во агентов>
```


## Structure

```
├── back-end
│ ├── agent - исходники агента
│ │ ├── api - сюда должна генерироваться документация
│ │ ├── cmd - точка входа
│ │ ├── internal
│ │ │ ├── agent_errors - кастомные внутренние ошибки
│ │ │ ├── app - инициализация и сборка приложения
│ │ │ ├── config - конфигурация агента
│ │ │ ├── services - слой сервиса
│ │ │ └── transport - слой транспорта (rest, kafka)
│ ├── db - папка с файлом БД для sqlite3
│ ├── models - общие модели сообщений и таблиц БД
│ └── orkestrator - исходники для оркестратора
│     ├── api - сюда должна складываться документация по api
│     ├── cmd - точка входа
│     ├── internal
│     │ ├── app - инициализация и сборка приложения
│     │ ├── clierrs - кастомные клиентские ошибки
│     │ ├── config - конфигурация
│     │ ├── interfaces - интерфейсы для соединения слоев
│     │ ├── repository - слой данных (хранилища)
│     │ ├── services - слой сервиса и бизнес логики
│     │ └── transport - транспортный слой (rest, kafka broker)
└── front-end - здесь должен быть рабочий фронт((
```