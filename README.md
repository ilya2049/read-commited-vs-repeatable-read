# Read commited vs repeatable read

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

Создаем и заполняем таблицу с пользователями

``` sql
create table users (
  id int primary key,
  name text not null,
  active bool not null
);

insert into users (id, name, active) values 
(1, 'Alex', true), 
(2, 'Sam', true), 
(3, 'Felix', true);
```

## Демонстрания аномалии _несогласованное чтение_

Запускаем _virus_

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/virus/main.go
```

Несколько раз запускаем _printer_, ловим баг

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/printer/main.go
```

## Сравниваем время выполнения только читающих транзакций

Создаем таблицу _passwords_ для тестов.

``` sql
create table passwords as
select
  generate_series(1, 10000) as id,
  substr(md5(random()::text), 0, 25) as hash;
```

Уровень изоляции Read Commited

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/txrace/main.go -mode=read-committed -reads=10
```

Уровень изоляции Repeatable Read

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/txrace/main.go -mode=repeatable-read -reads=10
```

Уровень изоляции Repeatable Read, режим доступа Read Only

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/txrace/main.go -mode=repeatable-read-read-only -reads=10
```