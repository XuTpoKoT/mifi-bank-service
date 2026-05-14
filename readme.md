# REST API банковского сервиса

## Реализовано
- Регистрация пользователей с проверкой уникальности.
- Аутентификация пользователей.
- Создание банковских счетов и управление ими.
- Операции с картами: генерация, просмотр, оплата.
- Переводы между счетами и пополнение баланса.
- Кредитные операции: оформление кредита, график платежей.

## Не реализовано
- Аналитика финансовых операций.
- Интеграция с ЦБ.
- SMTP — для отправки уведомлений по электронной почте.

## Запуск приложения

Выполните команду в корне проекта для поднятия бд:
```bash
docker-compose up -d
```

запустите приложение
```bash
go run cmd/server/main.go
```

## Тестирование функционала
**Базовые переменные**
```http
@host = http://localhost:8080
@token = {{token}}
```

### 1. Регистрация

#### Request
```http
POST {{host}}/register
Content-Type: application/json

{
  "username": "dima",
  "email": "dima@mail.com",
  "password": "dima"
}
```

#### Response
```json
{
  "id": 1,
  "username": "dima",
  "email": "dima@mail.com"
}
```

---

### 2. Логин

#### Request
```http
POST {{host}}/login
Content-Type: application/json

{
  "email": "dima@mail.com",
  "password": "dima"
}
```

#### Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Save token
```http
> {% client.global.set("token", response.body.token); %}
```

---

### 3. Создать счёт

#### Request
```http
POST {{host}}/accounts
Authorization: Bearer {{token}}
```

#### Response
```json
{
  "id": 1,
  "user_id": 1,
  "balance": 0,
  "currency": "RUB"
}
```

---

### 4. Пополнение счёта

#### Request
```http
POST {{host}}/accounts/topup
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "account_id": 1,
  "amount": 5000
}
```

#### Response
```json
{
  "message": "top up successful",
  "balance": 5000
}
```

---

### 5. Перевод

#### Request
```http
POST {{host}}/transfer
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 1000
}
```

#### Response
```json
{
  "message": "transfer completed",
  "transaction_id": 10
}
```

---

### 6. Выпуск карты

#### Request
```http
POST {{host}}/cards
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "account_id": 1
}
```

#### Response
```json
{
  "id": 1,
  "account_id": 1,
  "pan": "**** **** **** 1234",
  "expiry": "12/30"
}
```

---

### 7. Получить счета

#### Request
```http
GET {{host}}/accounts
Authorization: Bearer {{token}}
```

#### Response
```json
[
  {
    "id": 1,
    "balance": 5000,
    "currency": "RUB"
  }
]
```

---

### 8. Получить карты

#### Request
```http
GET {{host}}/cards
Authorization: Bearer {{token}}
```

#### Response
```json
[
  {
    "id": 1,
    "account_id": 1,
    "pan": "**** **** **** 1234",
    "expiry": "12/30"
  }
]
```

---

### 9. Кредит

#### Request
```http
POST {{host}}/credits
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "account_id": 1,
  "principal": 100000,
  "term_months": 12
}
```

#### Response
```json
{
  "message": "credit created"
}
```

---

### 10. График платежей

#### Request
```http
GET {{host}}/credits/1/schedule
Authorization: Bearer {{token}}
```

#### Response
```json
[
  {
    "id": 1,
    "credit_id": 1,
    "due_date": "2026-06-14",
    "amount": 9263.45,
    "status": "PENDING",
    "penalty": 0
  }
]
```