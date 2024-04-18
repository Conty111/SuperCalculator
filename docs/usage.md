# SuperCalculator

## Table of contents

- ### [Description](..%2FREADME.md#description)
- ### [How to run?](..%2FREADME.md#running)
- ### [How to use?](#how-to-use)
- ### [How it works?](FhowItWorks.md)
- ### [Project sructure](..%2FREADME.md#structure)

## How to use
## Requests

Если вы пользуетесь Postman, то можно импортировать коллекцию Calculator.postman_collection.json
Таймаут ответа от агента и период повторной отправки указываются в секундах, время выполнения операций - в миллисекундах

### Curl requests
Register
```
```
Create task
```
curl -L 'http://localhost:8000/api/v1/manager' \
-H 'Content-Type: application/json' \
-d '{
    "expression": "11+11+11"
}'
```
Get all tasks
```
curl -L 'http://localhost:8000/api/v1/manager/tasks'
```
Set settings (operation time in millisecond, timeout and time retry in seconds)
```
curl -L -X PUT 'http://localhost:8000/api/v1/manager/settings' \
-H 'Content-Type: application/json' \
-d '{
    "timeout_response": 10,
    "time_retry": 10,
    "add_time": 5000,
    "division_time": 6000,
    "subtract_time": 10,
    "multiply_time": 5
}'
```
Get workers info
```
curl -L 'http://localhost:8000/api/v1/manager/workers'
```