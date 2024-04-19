# SuperCalculator

## Table of contents

- ### [What is it?](#what-is-it)
- ### [How to run?](#how-to-run)
- ### [How to use?](docs%2Fusage.md#how-to-use)
- ### [How it works?](docs%2FhowItWorks.md)
- ### [Project sructure](#project-structure)

## What is it?

**SuperCalculator** - это распределенный калькулятор или проект, имитирующий отказоустойчивую
систему орекстрации трафика/задач/сообщений между сущностями. 

Также реализован пользовательский функционал (регистрация и авторизация, разграничение доступа на основе ролей, просмотр своих задач, администрирование пользователей и т.п.)

_P. S. front-end часть в проекте на данный отсутствует._

## How to run

### Requirements
Для запуска системы нужен **Go v1.21.6** (и выше) и **Docker v26.0.1** вместе с **Docker Compose**

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
2. Локально:
   1. Установить зависимости с помощью ```go mod tidy```
   2. Запуск Kafka
   ```
   export AGENTS_COUNT=<макс. кол-во агентов, которое будет использоваться>
   docker compose up kafka -d
   ```
   3. Запуск оркестратора
   ```
   make run-orkestrator
   ```
   4. Запуск агента(-ов) (в _<id_agenta>_ подставлять индекс агента из [system_config.json](system_config.json))
   ```
   go run -v ./back-end/agent/cmd/app/main.go serve <id_агента>
   ```




## Project structure

```
├── back-end
│ ├── agent - исходники агента
│ │ ├── cmd - точка входа
│ │ ├── internal
│ │ │ ├── agent_errors - кастомные внутренние ошибки
│ │ │ ├── app - инициализация и сборка приложения
│ │ │ ├── config - извлечение конфигурации
│ │ │ ├── services - слой сервиса
│ │ │ └── transport - транспортный слой (htpp, kafka, grpc)
│ ├── db - папка с файлами для БД
│ ├── models - модели сущностей приложения (сообщения, модели БД)
│ └── orkestrator - исходники для оркестратора
│     ├── cmd - точка входа
│     ├── internal
│     │ ├── app - инициализация и сборка приложения
│     │ ├── clierrs - кастомные клиентские ошибки
│     │ ├── config - извлечение конфигурации
│     │ ├── interfaces - интерфейсы для соединения слоев
│     │ ├── repository - слой данных (хранилища)
│     │ ├── services - слой сервиса и бизнес логики
│     │ └── transport - транспортный слой (htpp, kafka, grpc)
└── front-end - здесь должен был быть рабочий фронт((
```