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

## Эксперимент

Запускаем _virus_

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/virus/main.go
```

Несколько раз запускаем _printer_, ловим баг

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/printer/main.go
```