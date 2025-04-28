# backend

## Работа с миграциями

1. Установка golang-migrate
```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

2. Выполнение миграции
```shell
migrate -database "postgres://root:123@postgres:5432/domashka?sslmode=disable" -path ./migrations up
```

3. Откат миграции
```shell
migrate -database "postgres://root:123@postgres:5432/domashka?sslmode=disable" -path ./migrations down
```
