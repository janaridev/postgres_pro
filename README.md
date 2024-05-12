<h1>Postgres Pro Intern</h1>

Приложение, предоставляющее REST API для запуска команд - bash-скриптов

<h2>Требования для запуска приложения:</h2>
<ul>
    <li>Git</li>
    <li>Docker</li>
    <li>Docker-Compose</li>
</ul>

<h2>Инструкция по запуску приложения:</h2>

```bash
git clone git@github.com:janaridev/postgres_pro.git
```

```bash
cd postgres_pro
```

```bash
touch .env
```

Пример .env файла:

```bash
ENV=dev

# Postgres
PG_USER=admin
PG_PASSWORD=admin
PG_DB_NAME=postgres_pro
PG_EXTERNAL_PORT=5432
PG_INTERNAL_PORT=5432
PG_USE_SSL=disable
PG_HOST=postgres-pro

# PG Admin
PGADMIN_EMAIL=admin@gmail.com
PGADMIN_PASSWORD=admin
PGADMIN_EXTERNAL_PORT=8080
PGADMIN_INTERNAL_PORT=80

# Server
HOST=0.0.0.0
PORT=3000
```

Запуск:

```bash
make compose_up
```

<h2>Документация по апи:</h2>
<a href='https://documenter.getpostman.com/view/23019615/2sA3JNZzQG'>клик</a>

<h2>Дополнительные вещи которые были сделаны:</h2>
<ul>
    <li>Все задания со звездочкой были сделаны</li>
    <li>Использование docker && docker compose</li>
    <li>Дополнительно с ci был реализован push нового докер образа в docker hub</li>
    <li>Использование миграций с контейнером миграций docker compose</li>
    <li>Запуск апи на веб сервере (nginx)</li>
    <li>Админка для постгреса (pg admin)</li>
    <li>Эндпоинт для повторного запуска команды</li>
    <li>Использование триггеров для добавления/обновления полей created_at и updated_at</li>
</ul>
