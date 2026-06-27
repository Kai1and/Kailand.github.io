# Equipment Booking

Скелет приложения для бронирования оборудования.

## Структура

- `backend` - REST API на Go.
- `frontend` - интерфейс на React + Vite.
- `docs` - документация проекта.

Описание для курсовой находится в `course-description.md`.

## Быстрый старт

1. Создайте PostgreSQL базу `equipment_booking`.
2. Примените миграцию `backend/migrations/001_init.sql`.
3. Настройте `backend/.env` по примеру `backend/.env.example`.
4. Запустите API:

```bash
cd backend
go run ./cmd
```

5. Запустите frontend:

```bash
cd frontend
npm install
npm run dev
```
