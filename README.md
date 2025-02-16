# Avito-test

## Запуск проекта

1. Переименуйте файл `.env.example` в `.env`.
2. Запустите проект с помощью Docker Compose:
   ```sh
   docker-compose up --build
   ```

## API Эндпоинты

### 1. Аутентификация
**Запрос:**
```sh
curl -X POST "http://localhost:8080/api/auth" \
     -H "Content-Type: application/json" \
     -d '{"username": "test_user", "password": "password123"}'
```
**Ответ:**
```json
{
  "token": "your_jwt_token"
}
```

### 2. Получить информацию о монетах, инвентаре и истории транзакций
**Запрос:**
```sh
curl -X GET "http://localhost:8080/api/info" \
     -H "Authorization: Bearer your_jwt_token"
```

### 3. Отправить монеты другому пользователю
**Запрос:**
```sh
curl -X POST "http://localhost:8080/api/sendCoin" \
     -H "Authorization: Bearer your_jwt_token" \
     -H "Content-Type: application/json" \
     -d '{"toUser": "user2", "amount": 10}'
```

### 4. Купить предмет за монеты
**Запрос:**
```sh
curl -X GET "http://localhost:8080/api/buy/item_name" \
     -H "Authorization: Bearer your_jwt_token"
```

**Таблица мерча:**

| Название     | Цена |
|--------------|------|
| t-shirt      | 80   |
| cup          | 20   |
| book         | 50   |
| pen          | 10   |
| powerbank    | 200  |
| hoody        | 300  |
| umbrella     | 200  |
| socks        | 10   |
| wallet       | 50   |
| pink-hoody   | 500  |

## Используемые технологии
- **Язык программирования:** Go
- **Фреймворки и библиотеки:**
    - [ogen](https://github.com/ogen-go/ogen) - генерация API по OpenAPI спецификации
    - [squirrel](https://github.com/Masterminds/squirrel) - построение SQL-запросов
    - [go-transaction-manager](https://github.com/avito-tech/go-transaction-manager) - управление транзакциями
    - [zlog](https://github.com/rs/zerolog) - логирование
    - [lo](https://github.com/samber/lo) - утилиты для работы с данными
- **База данных:** PostgreSQL
- **Кэширование:** Redis

