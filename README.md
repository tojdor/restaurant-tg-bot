# Restaurant Telegram Bot Backend

Backend is intended to be called only by the Telegram bot service. It is not a public API for Telegram users.

## Authentication and authorization

Every endpoint requires both headers:

```http
X-Bot-Secret: <BOT_INTERNAL_SECRET>
X-Telegram-User-ID: <Telegram user ID from Update.Message.From.ID or CallbackQuery.From.ID>
```

`BOT_INTERNAL_SECRET` is an environment variable shared only by the bot and backend. The backend compares it in constant time, then loads the user and role from `users.telegram_user_id`. A missing/incorrect secret returns `401`; an unknown Telegram user or insufficient role returns `403`.

The bot must never take `X-Telegram-User-ID` from message text or from the client. It must use the ID supplied by Telegram in the received update. Do not expose this backend directly to the Internet; place it on a private network behind the bot service.

Roles:

- `admin` — menu, tables, and staff management.
- `waiter` — creates and manages only their own orders and order items.
- `kitchen` — reads orders and marks order items ready.

`telegram_user_id` is mandatory when creating a staff user. It is stable; Telegram usernames are not.

## Endpoints

All examples omit the two required authentication headers shown above.

| Method and path | Allowed role | What it does |
| --- | --- | --- |
| `POST /users` | admin | Creates a staff member associated with a Telegram user ID. |
| `GET /users` | admin | Lists staff members. |
| `DELETE /users/{id}` | admin | Removes a staff member. |
| `POST /dishes` | admin | Adds a dish to the menu. |
| `GET /dishes/category/{category}` | admin, waiter, kitchen | Lists dishes in a category. |
| `GET /dishes/name/{name}` | admin, waiter, kitchen | Returns a dish by its name. |
| `DELETE /dishes/{id}` | admin | Removes a dish from the menu. |
| `POST /tables` | admin | Creates a restaurant table. |
| `GET /tables` | admin, waiter, kitchen | Lists restaurant tables and their statuses. |
| `DELETE /tables/{number}` | admin | Removes a table. |
| `POST /orders` | waiter | Creates an order. The waiter is set from `X-Telegram-User-ID`; do not trust `waiter_id` in JSON. |
| `GET /orders` | admin, kitchen | Lists all orders. |
| `GET /orders/waiter/{waiter_id}` | admin, waiter | Lists a waiter's orders. A waiter may request only their own internal user ID. |
| `GET /orders/table/{table_number}` | admin, waiter, kitchen | Gets the order for a table. A waiter may view only their own order. |
| `PATCH /orders/{id}` | admin, waiter | Updates an order. A waiter may update only their own order and cannot reassign it. |
| `DELETE /orders/{id}` | admin, waiter | Deletes an order. A waiter may delete only their own order. |
| `POST /order-items` | waiter | Adds an item to the waiter's own order. |
| `GET /order-items/{id}` | admin, waiter, kitchen | Lists an order's items. A waiter may read only their own order. |
| `PATCH /order-items/{order_id}/{menu_item_id}/ready` | kitchen | Marks an order item ready or not ready. |
| `PATCH /order-items/{order_id}/{menu_item_id}/count` | waiter | Changes item quantity in the waiter's own order. |
| `DELETE /order-items/{order_id}/{menu_item_id}` | waiter | Deletes an item from the waiter's own order. |

## Request examples

Create a staff user (admin only):

```json
{
  "telegram_user_id": 123456789,
  "nickname": "Maria",
  "phone_number": "+79990000000",
  "role": "waiter"
}
```

Create a dish (admin only):

```json
{
  "name": "Caesar salad",
  "price": 490,
  "category": "salads"
}
```

Create an order (waiter only):

```json
{
  "table_number": 5,
  "is_served": false,
  "is_payed": false
}
```

Create an order item (waiter only):

```json
{
  "order_id": 42,
  "menu_item_id": 7,
  "count": 2,
  "is_ready": false
}
```

## Configuration

Required environment variables:

```dotenv
DATABASE_URL=postgres://user:password@host:5432/restaurant?sslmode=disable
BOT_INTERNAL_SECRET=use-a-long-random-secret-here
```

Before deploying, create the first admin directly in the database with a valid `telegram_user_id`, or prepare an administrative bootstrap migration. Never provide a public endpoint that lets an unknown Telegram user assign themselves the `admin` role.
