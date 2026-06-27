ALTER TABLE users ADD COLUMN IF NOT EXISTS ban_reason TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS ban_evidence TEXT NOT NULL DEFAULT '';

ALTER TABLE equipment ADD COLUMN IF NOT EXISTS hidden BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS available BOOLEAN NOT NULL DEFAULT true;

INSERT INTO users (name, email, password_hash, role, phone, city, avatar_url, bio)
VALUES
  ('Марина Орлова', 'marina@equipment.local', '$2a$10$zjVBxibw8De5HqJG5ftTr.lBT1sK/MoRJWr5h0oq99PA1QRIPqGvu', 'user', '+7 900 000-00-03', 'Екатеринбург', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&w=300&q=80', 'Помогаю с оборудованием для лекций и мероприятий.'),
  ('Илья Соколов', 'ilya@equipment.local', '$2a$10$zjVBxibw8De5HqJG5ftTr.lBT1sK/MoRJWr5h0oq99PA1QRIPqGvu', 'user', '+7 900 000-00-04', 'Челябинск', 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?auto=format&fit=crop&w=300&q=80', 'Сдаю ноутбуки, свет и звук для проектов.'),
  ('Диана Нуриева', 'diana@equipment.local', '$2a$10$zjVBxibw8De5HqJG5ftTr.lBT1sK/MoRJWr5h0oq99PA1QRIPqGvu', 'user', '+7 900 000-00-05', 'Пермь', 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=300&q=80', 'Камеры, микрофоны и свет для аккуратных съемок.')
ON CONFLICT (email) DO NOTHING;

INSERT INTO categories (name, description)
VALUES
  ('Звук', 'Микрофоны, рекордеры и акустика'),
  ('Свет', 'Осветительные приборы и стойки'),
  ('Инструменты', 'Техническое оборудование для монтажа и ремонта')
ON CONFLICT (name) DO NOTHING;

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'marina@equipment.local'), (SELECT id FROM categories WHERE name = 'Презентации'), 'Интерактивная панель 65 дюймов', 'Сенсорная панель для защиты проектов, лекций и командной работы.', 'PNL-4001', 'https://images.unsplash.com/photo-1593642634367-d91a135587b5?auto=format&fit=crop&w=1200&q=80', 'Екатеринбург', 2400, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'PNL-4001');

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'ilya@equipment.local'), (SELECT id FROM categories WHERE name = 'Компьютеры'), 'MacBook Pro 14 M2', 'Ноутбук для монтажа, дизайна и презентаций. Зарядка и чехол включены.', 'NBK-5001', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?auto=format&fit=crop&w=1200&q=80', 'Челябинск', 2100, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'NBK-5001');

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'diana@equipment.local'), (SELECT id FROM categories WHERE name = 'Звук'), 'Петличные микрофоны Rode Wireless GO II', 'Комплект из двух передатчиков для интервью, лекций и записи видео.', 'MIC-6001', 'https://images.unsplash.com/photo-1590602847861-f357a9332bbc?auto=format&fit=crop&w=1200&q=80', 'Пермь', 800, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'MIC-6001');

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'user@equipment.local'), (SELECT id FROM categories WHERE name = 'Свет'), 'Комплект света Godox SL60W', 'Два источника, стойки и софтбоксы для предметной и портретной съемки.', 'LGT-7001', 'https://images.unsplash.com/photo-1492691527719-9d1e07e534b4?auto=format&fit=crop&w=1200&q=80', 'Екатеринбург', 1100, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'LGT-7001');

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'marina@equipment.local'), (SELECT id FROM categories WHERE name = 'Звук'), 'Портативная колонка JBL PartyBox', 'Колонка для мероприятий, мастер-классов и небольших выступлений.', 'AUD-8001', 'https://images.unsplash.com/photo-1545454675-3531b543be5d?auto=format&fit=crop&w=1200&q=80', 'Екатеринбург', 900, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'AUD-8001');

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'ilya@equipment.local'), (SELECT id FROM categories WHERE name = 'Инструменты'), '3D-принтер Creality Ender 3 S1', 'Печать прототипов и учебных макетов, пластик PLA включен в стартовый комплект.', 'PRN-9001', 'https://images.unsplash.com/photo-1612637890088-8f9ba23b3f8f?auto=format&fit=crop&w=1200&q=80', 'Челябинск', 1300, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'PRN-9001');
