# API

Base URL: `http://localhost:8080/api`

## Auth

- `POST /auth/register` - регистрация пользователя.
- `POST /auth/login` - вход, возвращает JWT.
- `GET /auth/me` - текущий пользователь.

## Categories

- `GET /categories` - список категорий.
- `POST /categories` - создать категорию, admin.
- `PUT /categories/{id}` - обновить категорию, admin.
- `DELETE /categories/{id}` - удалить категорию, admin.

## Equipment

- `GET /equipment` - список оборудования.
- `GET /equipment/{id}` - карточка оборудования.
- `POST /equipment` - создать оборудование, admin.
- `PUT /equipment/{id}` - обновить оборудование, admin.
- `DELETE /equipment/{id}` - удалить оборудование, admin.

## Bookings

- `GET /bookings` - список бронирований.
- `POST /bookings` - создать бронирование.
- `PATCH /bookings/{id}/cancel` - отменить бронирование.
- `PATCH /bookings/{id}/status` - изменить статус, admin.

## Users

- `GET /users` - список пользователей, admin.
- `PATCH /users/{id}/role` - изменить роль, admin.
