# go_songs_lib_API
API для серсиса, имитирующего библиотеку различных песен и авторов (групп)

# Как запустить в докере?
## Требования:
- Установленный docker
- Docker Desktop (для Windows)
## 1 Собираем контейнеры: 
```shell
docker compose up -d --build
```
## 2 Чтобы посмотреть логи приложения:
```shell
docker compose logs app
```
## 3 Чтобы заполнить таблицы БД минимальными данными (опционально):
```shell
docker exec -it db_container pg_restore -U postgres -d songs_db /backups/songsData.sql
```

# Как запустить локально?
## Требования:
- Golang (1.23+)
- PostgreSQL 16.0
## 1 Устанавливаем зависимости:
```shell
go mod download
```
## 2 Нужно поменять host внутри DB_URL и POSTGRES_URL. Поменять SERV_PORT на 8080:
```shell
DB_URL='host=db_container...' --> DB_URL='host=localhost...'
POSTGRES_URL='host=db_container...' --> POSTGRES_URL='host=localhost...'
SERV_PORT=8081 --> SERV_PORT=8080
```
## 3 Запустить файл main.go:
```shell
go run main.go
```

# Как посмотреть документацию API?
## Вводим в браузерной строке:
```shell
http://localhost:8080/swagger/
```

## URL для сервера, выдающего расширенную информацию о песнях, необходимо указать в .env в INFO_API_URL
