# Read commited vs repeatable read

## Подготовка инфраструктуры

Запускаем контейнер с СУБД postgresql

``` sh
docker-compose up -d
```

Подключаемся к контейнеру c postgresql

``` sh
docker exec -it postgres bash
```

Подключаемся к СУБД postgresql из контейнера

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

Запускаем 'вредную' утилиту шутника

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/jocker/main.go
```

Несколько раз запускаем утилиту админа, ловим баг

``` sh
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres go run cmd/admin/main.go
```