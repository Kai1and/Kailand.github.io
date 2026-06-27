# Быстрое развертывание

## 1. База данных Neon

1. Создайте бесплатный проект PostgreSQL.
2. Откройте SQL Editor.
3. По порядку выполните файлы:
   - `backend/migrations/001_init.sql`
   - `backend/migrations/002_roles_and_demo_data.sql`
   - `backend/migrations/003_booking_chat_ban_demo.sql`
   - `backend/migrations/004_moderation_unread_attachments.sql`
4. Скопируйте строку подключения и сохраните ее как `DATABASE_URL`.

## 2. Backend на Render

Создайте Blueprint из репозитория: Render прочитает `render.yaml`.

Укажите:

- `DATABASE_URL` — строка подключения Neon;
- `ALLOWED_ORIGINS` — адрес frontend без завершающего `/`.

После запуска проверьте:

`https://<render-service>.onrender.com/api/health`

## 3. Frontend на Vercel

Создайте проект из того же репозитория:

- Root Directory: `frontend`
- Framework Preset: Vite
- Build Command: `npm run build`
- Output Directory: `dist`
- Environment Variable:
  `VITE_API_URL=https://<render-service>.onrender.com/api`

После первого deploy скопируйте адрес Vercel в `ALLOWED_ORIGINS` сервиса Render
и выполните повторный deploy backend.
