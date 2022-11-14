![poster](resourcer/poster.png)
# Use_balance_API

## Описание

Сервис, который позволяет работать с балансом пользователя

## Стек

- `Go`
- `Postgres`
- Фреймворк [Gin](https://github.com/gin-gonic/gin)
- `Docker`
- Конфигурация приложения [cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- Работа с БД [sqlx](https://github.com/jmoiron/sqlx)
- Логер [logrus](https://github.com/sirupsen/logrus)

## Реализованный функционал

- [x] Основной функционал работы с пользователем
- [x] Покрыт код тестами
- [x] Swagger [swagger](http://localhost:8080/swagger/index.html#)
- [x] Дополнительный функционал (Доп 1, Доп 2)

## Запуск

1. Склонировать репозиторий

```
git clone https://github.com/vladjong/user_balance.git
```

2. Включить Docker

3. Открыть терминал и набрать:
```
make
```
По стандарту запускается цель `docker-compose`

4. Тесты

```
make test
```

5. Проверка на стиль

```
make lint
```

## Тестирование

### Стандартный порт `8080`

- `/api` - REST API

- `/swagger` - Swagger API

### Post

- `/:id/:val` Метод пополнения баланса пользователя

Curl:
```
curl -X 'POST' \
  'http://localhost:8080/api/1/1000' \
  -H 'accept: application/json' \
  -d ''
```
Response body:
```
{
  "Status": "ok"
}
```

- `/reserv/:id/:id_ser/:id_ord/:val` Метод резервирования средств с основного баланса на отдельном счете

Curl:
```
curl -X 'POST' \
  'http://localhost:8080/api/reserv/1/1/1/125' \
  -H 'accept: application/json' \
  -d ''
```
Response body:
```
{
  "Status": "ok"
}
```

- `/accept/:id/:id_ser/:id_ord/:val` Метод признания выручки - списывает из резерва деньги, добавляет данные в отче для бухгалтерии

Curl:
```
curl -X 'POST' \
  'http://localhost:8080/api/accept/1/1/1/125' \
  -H 'accept: application/json' \
  -d ''
```
Response body:
```
{
  "Status": "ok"
}
```

- `/accept/:id/:id_ser/:id_ord/:val` Метод разрезервирования денег - переводятся обратно на счет пользователя

Curl:
```
curl -X 'POST' \
  'http://localhost:8080/api/reject/1/1/1/125' \
  -H 'accept: application/json' \
  -d ''
```
Response body:
```
{
  "Status": "ok"
}
```

### Get

- `/:id` Метод получения баланса пользователя

Curl:
```
curl -X 'GET' \
  'http://localhost:8080/api/1' \
  -H 'accept: application/json'
```
Response body:
```
{
  "id": 1,
  "balance": "1000"
}
```

- `/report/:date` Метод получения месячного отчета

Curl:
```
curl -X 'GET' \
  'http://localhost:8080/api/report/2022-11' \
  -H 'accept: application/json'
```
Response body:
```
{
  "Filename": "data/report_2022-11.csv"
}
```

Пример отчета находится в `data/report_2022-11.csv`

| id | name         | all_sum |
|----|--------------|---------|
| 1  | Доставка     | 554.23  |
| 2  | Консультация | 845.83  |
| 3  | Пополнение   | 3000    |
| 4  | Упаковка     | 325     |

- `/history/:id/:dat` Метод получения месячного отчета для пользователя

Curl:
```
curl -X 'GET' \
  'http://localhost:8080/api/history/1/2022-11' \
  -H 'accept: application/json'
```

Response body:
```
[
  {
    "id": 1,
    "service_name": "Консультация",
    "order_name": "А3",
    "sum": "500",
    "status_transaction": true,
    "date": "2022-11-14T13:06:45.358993Z"
  },
  {
    "id": 2,
    "service_name": "Упаковка",
    "order_name": "А1",
    "sum": "250",
    "status_transaction": false,
    "date": "2022-11-14T13:06:08.14635Z"
  },
  {
    "id": 3,
    "service_name": "Доставка",
    "order_name": "А2",
    "sum": "500",
    "status_transaction": true,
    "date": "2022-11-14T13:05:52.131081Z"
  },
]
```

### Кейс 1: Совершение транзакции на сумму большей чем баланс клиента

Curl:
```
curl -X 'POST' \
  'http://localhost:8080/api/reserv/1/1/1/10000' \
  -H 'accept: application/json' \
  -d ''
```
Response body:
```
{
  "message": "error: customer balance less than transaction cost"
}
```

### Кейс 2: Одобрение не существующей транзакции клиента

Curl:
```
curl -X 'POST' \
  'http://localhost:8080/api/accept/1/1/1/10' \
  -H 'accept: application/json' \
  -d ''
```
Response body:
```
{
  "message": "error: this id don't exist"
}
```

### Кейс 3: Узнать баланс не существующего клиента

Curl:
```
curl -X 'GET' \
  'http://localhost:8080/api/5' \
  -H 'accept: application/json'
```
Response body:
```
{
  "message": "error: id don't exist"
}
```

### Кейс 4: Экспорт отчета в котором не было транзакции в определенный период

Curl:
```
curl -X 'GET' \
  'http://localhost:8080/api/report/2021-01' \
  -H 'accept: application/json'
```
Response body:
```
{
  "message": "don't have history report in 2021-01-01 00:00:00 +0000 UTC"
}
```

### Кейс 5: История клиента в котором не было транзакции в определенный период

Curl:
```
curl -X 'GET' \
  'http://localhost:8080/api/history/1/2021-01' \
  -H 'accept: application/json'
```
Response body:
```
{
  "message": "don't have customer id: 1 history report in 2021-01-01 00:00:00 +0000 UTC"
}
```

## Диаграмма БД

![db](resourcer/db.png)

### Таблица Customers
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор клиента       | id                 | |
| Баланс клиента                     | balance            | Актуальный баланс клиента |

### Таблица Accounts
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор промежуточного счета       | id                 | |
| Идентификатор клиента       | id                 | |
| Баланс клиента                     | balance            | Актуальный баланс клиента на промежуточном счете |

### Таблица Services
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор услуги       | id                 | |
| Имя услуги                 | name               | |

### Таблица Orders
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор заказа       | id                 | |
| Имя заказа                 | name               | |

### Таблица Transaction
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор транзакции       | id | |
| Идентификатор клиента                 | customer_id | |
| Идентификатор услуги                 | service_id | |
| Идентификатор заказа                 | order_id | |
| Сумма транзакции                 | cost | Сумма, которая перевелась на промежуточный счет |
| Дата транзакции                | transaction_datetime | Дата совершения транзакции |

### Таблица History
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор истории транзакций       | id | |
| Идентификатор транзакции   | transaction_id | |
| Дата применения операции    | accounitng_datetime | Дата списания денег с промежуточного счета |
| Статус транзакции     | status_transaction | true - успешно; false -отмененная |

### Таблица Expected_transaction
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор ожидаемой транзакции      | id | |
| Идентификатор транзакции   | transaction_id | |

### Представление History_report
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор истории транзакций       | id | |
| Название услуги   | name | |
| Cумма   | cost | |
| Дата применение транзакции   | accounting_datetime | Время пременения транзакции |

### Представление Customer_report
| **Поле**                    | **Название поля в системе** | **Описание**
|:---------------------------:|:---------------------------:|:------------:|
| Идентификатор истории транзакций  | id | |
| Название услуги   | service_name | |
| Название заказа   | order_name |  |
| Cумма   | sum | Сумма транзакции |
| Статус транзакции   | status_transaction | Время пременения транзакции |
| Дата применение транзакции   | date | |

#### Для тестирования таблицы Services и Orderes заполняются тестовыми данными

### Таблица Services
| **id**  | **name** |
|:-------:|:--------:|
| 1      | Упаковка       |
| 2      | Доставка     |
| 3      | Консультация     |
| 4      | Пополнение     |

### Таблица Orders
| **id**  | **name** |
|:-------:|:--------:|
| 1      | A1       |
| 2      | A2     |
| 3      | A3     |
| 4      | Баланс     |