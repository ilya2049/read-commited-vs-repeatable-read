# Сравнение времени выполнения только читающих транзакций

## Подготовка инфраструктуры

Запускаем контейнер с postgresql

``` sh
docker-compose up -d
```

Подключаемся к контейнеру c postgresql

``` sh
docker exec -it postgres bash
```

Подключаемся к postgresql из контейнера

``` sh
psql -U postgres
```

Создаем таблицу _passwords_.

``` sql
create table passwords as
select
  generate_series(1, 10000) as id,
  substr(md5(random()::text), 0, 25) as hash;
```

Эксперимент 1. Уровень изоляции Read Commited

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/txrace/main.go -mode=read-committed -reads=10
```

Эксперимент 2. Уровень изоляции Repeatable Read

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/txrace/main.go -mode=repeatable-read -reads=10
```

Эксперимент 3. Уровень изоляции Repeatable Read, режим доступа Read Only

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/txrace/main.go -mode=repeatable-read-read-only -reads=10
```