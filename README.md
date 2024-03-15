# content-api

API сервис для IOS-приложения content

## Развертывание

Создайте файл ```contentApi/.env``` с данными для базы данных PostgreSQL:
```
POSTGRES_NAME=
POSTGRES_PASSWORD=
POSTGRES_USER=
POSTGRES_HOST=db
POSTGRES_PORT=5432

TOKEN_HOUR_LIFESPAN=
API_SECRET=
```

Запустите докер:
```
docker compose up
```
