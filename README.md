# Restaurant Telegram Bot Backend

Backend предназначен **только для вызова сервисом Telegram-бота**.
Это не публичный API для пользователей Telegram.

---

## Authentication and Authorization

Каждый endpoint требует наличия **обоих** заголовков:

```http
X-Bot-Secret: <BOT_INTERNAL_SECRET>
X-Telegram-User-ID: <Telegram user ID из Update.Message.From.ID или CallbackQuery.From.ID>
```

### Как работает авторизация

1. Backend получает `X-Bot-Secret`.
2. Сравнивает его с `BOT_INTERNAL_SECRET` из переменных окружения.
3. Секрет сравнивается с использованием **constant-time comparison**.
4. Backend получает `X-Telegram-User-ID`.
5. По этому ID ищет пользователя в `users.telegram_user_id`.
6. После этого проверяет его роль.
7. Если роль разрешает выполнение операции — запрос выполняется.

### Коды ошибок

* `401 Unauthorized` — отсутствует или указан неправильный `X-Bot-Secret`.
* `403 Forbidden` — Telegram-пользователь не найден или у него недостаточно прав.

> **Важно:** бот никогда не должен брать `X-Telegram-User-ID` из текста сообщения или от пользователя.
>
> ID должен браться непосредственно из Telegram Update:
>
> * `Update.Message.From.ID`
> * `CallbackQuery.From.ID`

Backend не следует открывать напрямую в Интернет. Рекомендуется разместить его в **приватной сети** и разрешить доступ только сервису Telegram-бота.

---

## Roles

### `admin`

Администратор может:

* управлять персоналом;
* добавлять и удалять блюда;
* создавать и удалять столики;
* просматривать и изменять заказы.

### `waiter`

Официант может:

* просматривать меню;
* просматривать столики;
* создавать свои заказы;
* изменять свои заказы;
* добавлять позиции в свои заказы;
* изменять количество позиций;
* удалять позиции;
* просматривать только свои заказы.

Официант **не может**:

* изменять чужие заказы;
* переназначать заказ другому официанту;
* управлять персоналом;
* управлять меню;
* управлять столиками.

### `kitchen`

Кухня может:

* просматривать меню;
* просматривать столики;
* просматривать все заказы;
* просматривать позиции заказов;
* отмечать позиции как готовые или неготовые.

---

## Telegram User ID

При создании сотрудника `telegram_user_id` является обязательным:

```json
{
  "telegram_user_id": 123456789,
  "nickname": "Maria",
  "phone_number": "+79990000000",
  "role": "waiter"
}
```

`telegram_user_id` является стабильным идентификатором Telegram-пользователя.

**Не следует использовать Telegram username в качестве идентификатора пользователя**, поскольку username может отсутствовать или измениться.

---

# API Endpoints

Во всех примерах ниже обязательные заголовки аутентификации не показаны.

| Method   | Endpoint                                       | Role                         | Description                     |
| -------- | ---------------------------------------------- | ---------------------------- | ------------------------------- |
| `POST`   | `/users`                                       | `admin`                      | Создаёт сотрудника              |
| `GET`    | `/users`                                       | `admin`                      | Возвращает список сотрудников   |
| `DELETE` | `/users/{id}`                                  | `admin`                      | Удаляет сотрудника              |
| `POST`   | `/dishes`                                      | `admin`                      | Добавляет блюдо                 |
| `GET`    | `/dishes/category/{category}`                  | `admin`, `waiter`, `kitchen` | Возвращает блюда категории      |
| `GET`    | `/dishes/name/{name}`                          | `admin`, `waiter`, `kitchen` | Возвращает блюдо по названию    |
| `DELETE` | `/dishes/{id}`                                 | `admin`                      | Удаляет блюдо                   |
| `POST`   | `/tables`                                      | `admin`                      | Создаёт столик                  |
| `GET`    | `/tables`                                      | `admin`, `waiter`, `kitchen` | Возвращает столики и их статусы |
| `DELETE` | `/tables/{number}`                             | `admin`                      | Удаляет столик                  |
| `POST`   | `/orders`                                      | `waiter`                     | Создаёт заказ                   |
| `GET`    | `/orders`                                      | `admin`, `kitchen`           | Возвращает все заказы           |
| `GET`    | `/orders/waiter/{waiter_id}`                   | `admin`, `waiter`            | Возвращает заказы официанта     |
| `GET`    | `/orders/table/{table_number}`                 | `admin`, `waiter`, `kitchen` | Возвращает заказ столика        |
| `PATCH`  | `/orders/{id}`                                 | `admin`, `waiter`            | Обновляет заказ                 |
| `DELETE` | `/orders/{id}`                                 | `admin`, `waiter`            | Удаляет заказ                   |
| `POST`   | `/order-items`                                 | `waiter`                     | Добавляет позицию в заказ       |
| `GET`    | `/order-items/{id}`                            | `admin`, `waiter`, `kitchen` | Возвращает позиции заказа       |
| `PATCH`  | `/order-items/{order_id}/{menu_item_id}/ready` | `kitchen`                    | Изменяет статус готовности      |
| `PATCH`  | `/order-items/{order_id}/{menu_item_id}/count` | `waiter`                     | Изменяет количество             |
| `DELETE` | `/order-items/{order_id}/{menu_item_id}`       | `waiter`                     | Удаляет позицию                 |

---

# Request Examples

## Create Staff User

**Role:** `admin`

```http
POST /users
X-Bot-Secret: <BOT_INTERNAL_SECRET>
X-Telegram-User-ID: 123456789
Content-Type: application/json
```

```json
{
  "telegram_user_id": 123456789,
  "nickname": "Maria",
  "phone_number": "+79990000000",
  "role": "waiter"
}
```

---

## Create Dish

**Role:** `admin`

```http
POST /dishes
X-Bot-Secret: <BOT_INTERNAL_SECRET>
X-Telegram-User-ID: 123456789
Content-Type: application/json
```

```json
{
  "name": "Caesar salad",
  "price": 490,
  "category": "salads"
}
```

---

## Create Order

**Role:** `waiter`

```http
POST /orders
X-Bot-Secret: <BOT_INTERNAL_SECRET>
X-Telegram-User-ID: 123456789
Content-Type: application/json
```

```json
{
  "table_number": 5,
  "is_served": false,
  "is_payed": false
}
```

### Important

`waiter_id` **не должен передаваться в JSON**.

Backend должен определить официанта самостоятельно по:

```http
X-Telegram-User-ID
```

После этого backend находит соответствующего пользователя в базе данных и использует его внутренний `id`.

---

## Create Order Item

**Role:** `waiter`

```http
POST /order-items
X-Bot-Secret: <BOT_INTERNAL_SECRET>
X-Telegram-User-ID: 123456789
Content-Type: application/json
```

```json
{
  "order_id": 42,
  "menu_item_id": 7,
  "count": 2,
  "is_ready": false
}
```

Backend должен дополнительно проверить, что `order_id` принадлежит официанту, который отправил запрос.

---

# Configuration

Необходимые переменные окружения:

```dotenv
DATABASE_URL=postgres://user:password@host:5432/restaurant?sslmode=disable
BOT_INTERNAL_SECRET=use-a-long-random-secret-here
```

### `DATABASE_URL`

URL подключения к PostgreSQL:

```text
postgres://USER:PASSWORD@HOST:PORT/DATABASE?sslmode=disable
```

Например:

```dotenv
DATABASE_URL=postgres://postgres:password@localhost:5432/restaurant?sslmode=disable
```

### `BOT_INTERNAL_SECRET`

Секрет, который используется для проверки того, что запрос действительно пришёл от Telegram-бота.

Используйте длинную случайную строку.

Например:

```dotenv
BOT_INTERNAL_SECRET=very-long-random-secret
```

В production необходимо использовать действительно случайный секрет и **не хранить его непосредственно в исходном коде или Git-репозитории**.

---

# Initial Admin

Перед запуском production-системы необходимо создать первого администратора.

Это можно сделать:

1. непосредственно в базе данных;
2. через bootstrap migration;
3. через отдельный безопасный административный механизм.

Первый администратор должен иметь действительный:

```text
telegram_user_id
```

Например:

```sql
INSERT INTO users (
    telegram_user_id,
    nickname,
    phone_number,
    role
)
VALUES (
    123456789,
    'Admin',
    '+79990000000',
    'admin'
);
```

> **Важно:** никогда не создавайте публичный endpoint, через который неизвестный Telegram-пользователь может самостоятельно назначить себе роль `admin`.

---

# Security Requirements

Backend должен соблюдать следующие правила:

### 1. Никогда не доверять ID пользователя из JSON

Неправильно:

```json
{
  "waiter_id": 15,
  "table_number": 5
}
```

если `waiter_id` используется для определения владельца заказа.

Правильно:

```http
X-Telegram-User-ID: 123456789
```

Backend самостоятельно находит пользователя:

```text
Telegram User ID
        ↓
users.telegram_user_id
        ↓
internal users.id
        ↓
role
        ↓
authorization
```

### 2. Проверять владельца заказа

Если пользователь имеет роль `waiter`, он может работать только со своими заказами.

Например:

```text
waiter Telegram ID
        ↓
user.id = 15
        ↓
order.waiter_id = 15
```

Если:

```text
order.waiter_id != user.id
```

запрос должен быть отклонён с:

```http
403 Forbidden
```

### 3. Проверять роль

Каждый endpoint должен проверять, имеет ли пользователь необходимую роль.

Например:

```text
POST /dishes
        ↓
authentication
        ↓
load user
        ↓
check role
        ↓
admin?
   ┌────┴────┐
  yes       no
   ↓         ↓
 handler    403
```

### 4. Не открывать backend в Интернет

Backend предназначен для взаимодействия:

```text
Telegram
    ↓
Bot Service
    ↓
Private Network
    ↓
Restaurant Backend
    ↓
PostgreSQL
```

Пользователь Telegram не должен иметь прямого доступа к backend.

---

# Recommended Request Flow

Для каждого запроса рекомендуется следующая последовательность:

```text
HTTP Request
     ↓
X-Bot-Secret validation
     ↓
X-Telegram-User-ID validation
     ↓
Find user in database
     ↓
Check user exists
     ↓
Load user role
     ↓
Check permissions
     ↓
Handler
     ↓
Service
     ↓
Storage
     ↓
PostgreSQL
```

Таким образом, handlers не должны самостоятельно реализовывать всю логику авторизации, проверки ролей и работы с базой данных.

---

# Project Architecture

Рекомендуемое разделение:

```text
Handler
   ↓
Middleware
   ↓
Service
   ↓
Storage
   ↓
PostgreSQL
```

### Middleware

Отвечает за:

* проверку `X-Bot-Secret`;
* получение `X-Telegram-User-ID`;
* поиск пользователя;
* проверку авторизации;
* передачу пользователя/его роли дальше в request context.

### Handler

Отвечает за:

* HTTP;
* чтение JSON;
* валидацию HTTP-запроса;
* формирование HTTP-ответа.

### Service

Отвечает за:

* бизнес-логику;
* проверки владельца заказа;
* проверки разрешений, если они относятся к бизнес-логике;
* взаимодействие между несколькими storage-операциями.

### Storage

Отвечает только за:

* SQL;
* получение данных;
* создание данных;
* изменение данных;
* удаление данных.

---

# Summary

Backend предоставляет API для управления рестораном через Telegram-бота.

Основные сущности:

```text
Users
  ↓
Roles
  ↓
Tables
  ↓
Dishes
  ↓
Orders
  ↓
Order Items
```

Основные роли:

```text
admin
waiter
kitchen
```

Аутентификация строится на:

```text
BOT_INTERNAL_SECRET
+
X-Telegram-User-ID
```

Авторизация строится на:

```text
Telegram User ID
        ↓
Database User
        ↓
Role
        ↓
Permission
```

Backend не является публичным API и должен быть доступен только сервису Telegram-бота.
