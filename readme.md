# Project

## TODO
- кредиты 
- график платежей
- scheduler

## TODO
- аналитика
- SOAP ЦБ
- email

### регистрация
```bash
curl -X POST localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{
     "username":"dima",
     "email":"dima@mail.com",
     "password":"dima"
    }'
```

### login
```bash
curl -X POST localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{
     "email":"dima@mail.com",
     "password":"dima"
    }'
```


### создать счёт
```bash
curl -X POST localhost:8080/accounts \
-H "Authorization: Bearer TOKEN"
```

### пополнить
```bash
curl -X POST localhost:8080/accounts/topup \
    -H "Authorization: Bearer TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "account_id":1,
      "amount":5000
    }'
```

### перевод
```bash
curl -X POST localhost:8080/transfer \
    -H "Authorization: Bearer TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "from_account_id":1,
      "to_account_id":2,
      "amount":1000
    }'
```


