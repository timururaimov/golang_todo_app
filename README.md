# Golang ToDo API

REST API для управления пользователями и их задачами. Проект написан на Go, использует PostgreSQL и предоставляет интерактивную документацию Swagger.

## Возможности

- создание, просмотр, изменение и удаление пользователей;
- создание, просмотр, изменение и удаление задач;
- фильтрация задач по автору и постраничная выдача;
- статистика по задачам: число созданных и выполненных, доля выполнения и среднее время выполнения;
- версионированный API (`/api/v1`), graceful shutdown, middleware для request ID, логирования, CORS и обработки паник.

## Технологии

- Go 1.26.1+
- PostgreSQL
- `pgx` для работы с базой данных
- Docker Compose для локального окружения
- Swagger для документации API

## Быстрый запуск

Нужны Go 1.26.1 или новее, Docker и Docker Compose.

1. Скопируйте пример конфигурации и задайте параметры PostgreSQL:

   ```sh
   cp .env.example .env
   ```

   Заполните в `.env` значения `POSTGRES_USER`, `POSTGRES_PASSWORD` и `POSTGRES_DB`. Остальные настройки можно оставить по умолчанию. Файл `.env` содержит локальные секреты — не публикуйте его.

2. Запустите PostgreSQL:

   ```sh
   make env-up
   ```

3. Примените миграции:

   ```sh
   make migrate-up
   ```

4. Запустите API:

   ```sh
   make todoapp-run
   ```

Сервер будет доступен по адресу `http://localhost:8080`. В текущей конфигурации API подключается к базе на `localhost` и порту `POSTGRES_FORWARD_PORT` (по умолчанию `5432`); при необходимости запустите проброс порта командой `make env-port-forward`.

### Запуск приложения через Docker Compose

После настройки `.env` запустите приложение и базу:

```sh
make todoapp-deploy
```

Compose запускает API и PostgreSQL, но миграции нужно применить отдельно командой `make migrate-up`. Остановка приложения:

```sh
make todoapp-stop
```

## API

Базовый путь: `http://localhost:8080/api/v1`.

| Метод | Путь | Описание |
| --- | --- | --- |
| `POST` | `/users` | Создать пользователя |
| `GET` | `/users` | Получить пользователей; поддерживаются `limit` и `offset` |
| `GET` | `/users/{id}` | Получить пользователя |
| `PATCH` | `/users/{id}` | Изменить пользователя |
| `DELETE` | `/users/{id}` | Удалить пользователя |
| `POST` | `/tasks` | Создать задачу для пользователя |
| `GET` | `/tasks` | Получить задачи; поддерживаются `author_user_id`, `limit` и `offset` |
| `GET` | `/tasks/{id}` | Получить задачу |
| `PATCH` | `/tasks/{id}` | Изменить задачу |
| `DELETE` | `/tasks/{id}` | Удалить задачу |
| `GET` | `/statistics` | Получить статистику; поддерживаются `author_user_id`, `from`, `to` |

Пример создания пользователя:

```sh
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"Иван Иванов","phone_number":"+77001234567"}'
```

Пример создания задачи (подставьте существующий `author_user_id`):

```sh
curl -X POST http://localhost:8080/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Прочитать документацию","description":"Начать с README","author_user_id":1}'
```

Период статистики задаётся датами `from` и `to` в формате `YYYY-MM-DD`; укажите обе даты. Начало периода включается, конечная дата не включается. Например: `/api/v1/statistics?from=2026-01-01&to=2026-02-01`.

Интерактивная документация: [`http://localhost:8080/swagger/`](http://localhost:8080/swagger/).

## Конфигурация

Приложение читает параметры из переменных окружения (локально их удобно задать через `.env`):

| Переменная | Назначение | Значение по умолчанию |
| --- | --- | --- |
| `HTTP_ADDR` | Адрес HTTP-сервера | `:` (обычно используйте `:8080`) |
| `HTTP_SHUTDOWN_TIMEOUT` | Таймаут graceful shutdown | `30s` |
| `POSTGRES_HOST` | Хост PostgreSQL | обязательно задать |
| `POSTGRES_PORT` | Порт PostgreSQL | `5432` |
| `POSTGRES_USER` | Пользователь PostgreSQL | обязательно задать |
| `POSTGRES_PASSWORD` | Пароль PostgreSQL | обязательно задать |
| `POSTGRES_DB` | Имя базы данных | обязательно задать |
| `POSTGRES_TIMEOUT` | Таймаут подключения | `10s` |
| `LOGGER_LEVEL` | Уровень логирования | `DEBUG` |
| `LOGGER_FOLDER` | Каталог файлов логов | обязательно задать |
| `TIME_ZONE` | Часовой пояс | `UTC` |
| `CORS_ALLOWED_ORIGINS` | Разрешённые origin через запятую | пустой список |

При ручном запуске через `make todoapp-run` Makefile устанавливает каталог логов в `out/logs`, хост PostgreSQL в `localhost` и использует порт из `POSTGRES_FORWARD_PORT` (или `5432`).

Укажите в `CORS_ALLOWED_ORIGINS` origin веб-клиента, например `https://example.com`. Для нескольких адресов используйте список через запятую. В origin указываются схема, домен и при необходимости порт, без пути. Пустой список запрещает браузерные cross-origin запросы.

## Структура проекта

```text
cmd/todoapp/       точка входа и Dockerfile приложения
internal/core/     конфигурация, доменные типы, PostgreSQL и HTTP-инфраструктура
internal/features/ логика пользователей, задач и статистики
migrations/        SQL-миграции базы данных
docs/              сгенерированная спецификация Swagger
```

