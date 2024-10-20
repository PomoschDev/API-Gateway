# Маршрутизатор endpoints

## Конфигурация
Конфигурация находится в ```./config/prod.yaml```
```yaml
env: "prod" #Прод конфиг
api_server: #Настройки для api сервера
  port: 8010 #Порт который будет прослушивать сервер
  timeout: 5s #Таймаут запроса
grpc_server: #Настройки для gprc сервера
  host: localhost #Хост на котором находится grpc сервис с базой данных
  port: 44044 #Порт который прослушивает grpc сервис
  timeout: 5s #Таймаут запроса
swagger: false #Запускать ли сваггер-документацию (true - включить)
```

## Запуск
Присутствует запуск с аргументами: ```apigateway --config=./config/prod.yaml```

## Swagger
Swagger-документация запускается вместе с основным сервером, по умолчанию документация находится по адресу 
```http://localhost:8010/swagger/```. Для того что бы документация была доступна, в файле конфигурации параметр 
**swagger** нужно установить в **true**, по умолчанию в ```prod.yaml``` данная опция выключена.

## Docker
В корне присутствует **Dockerfile** и **docker-compose.yml**. В случае клонирования репозитория, достаточно
запустить команду ```docker compose up -d``` для разворачивания проекта. По умолчанию будет скопирована конфигурация
из папки ```./config/prod.yaml``` и скопирована в контейнер как ```local.yaml```, сервис запустится с аргументом
```--config=local.yaml```, порт контейнера по умолчанию **8010**
### Docker без клонирования репозитория
В папке **```docker```**, так же присутствует **Dockerfile** и **docker-compose.yml**, данные файлы предназначены
для того, что бы не клонировать весь репозиторий, сервер склонируется прямиков в контейнере, рядом с docker-файлами
обязательно должен лежать файл конфигурации ```prod.yaml```, в следствии он будет скопирован в контейнере под именем
```local.yaml```
для
запуска сервиса с аргументом ```--config=local.yaml```. **Рекомендуется использовать именно этот вариант
развертывания сервера, так при развертывании мы всегда будем получать актуальную версию сервиса.**