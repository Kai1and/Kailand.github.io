# Database

Проект использует PostgreSQL.

## Основные таблицы

- `users` - пользователи и роли.
- `categories` - категории оборудования.
- `equipment` - единицы оборудования.
- `bookings` - заявки на бронирование.

## Статусы бронирования

- `pending`
- `approved`
- `rejected`
- `cancelled`
- `returned`

## Миграции

Начальная схема находится в `backend/migrations/001_init.sql`.
