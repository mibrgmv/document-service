## Инструменты
Go 1.24, Gin, PostgreSQL, pgx, Redis, JWT, golang-migrate, Swagger, testify, mock

### Что есть хорошего 
- Генерация swagger документации
- Кэш в Redis с инвалидацией
- JWT авторизация
- Контейнеризация в docker  
- Юнит тесты для слоя сервисов
- Конфигурация из `yaml` и `env` файлов

### Что можно добавить
- Метрики и логирование
- Разделить сервис авторизации и сервис документов на разные микросервисы
- Авторизация через identity provider (например Keycloak)
- Пагинация для `GET /api/docs`

## Инструкция для запуска
- склонировать репозиторий
- локально
  - создать `.env` и заполнить его значениями 
  - установить зависимости `go mod download`
  - поднять окружение `docker-compose up -d postgres redis`
  - запустить `go run cmd/app/main.go`
- через Docker
  - `docker compose up -d`
- наслаждаться на порту `http://localhost:8080`, swagger на `http://localhost:8080/swagger`
- **note:** миграции поднимутся автоматически

### Пример `.env`
```dotenv
SERVER_PORT=8080
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=documents
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=56dhu8ytvf
ADMIN_TOKEN=f86jno7rcbu
```

## Описание API
- `POST /api/register` - регистрация нового пользователя
- `POST /api/auth` - аутентификация, получение JWT токена
- `DELETE /api/auth/{token}` - завершение сессии
- `GET /api/docs` - список документов с фильтрацией
- `POST /api/docs` - загрузка нового документа
- `GET /api/docs/{id}` - получение документа по ID
- `DELETE /api/docs/{id}` - удаление документа

## Тестирование через Swagger
1. Регистрация
```
POST /api/register
{
  "token": "admin-token",
  "login": "test",
  "pswd": "Password123!"
}
```
2. Авторизация
```
POST /api/auth
{
  "login": "test",
  "pswd": "Password123!"
}
```
3. Загрузка файла
- `POST /api/docs`
- Authorization: `Bearer TOKEN`
- Form-data
    - `meta`: `{"name":"test.txt","file":true,"public":true,"mime":"text/plain","grant":[]}`
    - `file`: выбрать файл
4. Загрузка JSON'а
- `POST /api/docs`
- Authorization: `Bearer TOKEN`
- Form-data
  - `meta`: `{"name":"data.json","file":false,"public":true,"mime":"application/json","grant":[]}`
  - `json`: `{"key": "value", "number": 123}`
5. Получение списка документов
- `GET /api/docs`
- Authorization: `Bearer TOKEN`
6. Получение документа по ID
- `GET /api/docs/{id}`
- Authorization: `Bearer TOKEN`
