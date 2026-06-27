ALTER TABLE users ADD COLUMN IF NOT EXISTS blocked BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS hidden BOOLEAN NOT NULL DEFAULT false;

DELETE FROM equipment current
USING equipment duplicate
WHERE current.serial = duplicate.serial AND current.id > duplicate.id;

CREATE UNIQUE INDEX IF NOT EXISTS equipment_serial_unique_idx ON equipment(serial) WHERE serial <> '';

INSERT INTO users (name, email, password_hash, role, phone, city, avatar_url, bio)
VALUES
  ('Администратор', 'admin@equipment.local', '$2a$10$zjVBxibw8De5HqJG5ftTr.lBT1sK/MoRJWr5h0oq99PA1QRIPqGvu', 'admin', '+7 900 000-00-01', 'Екатеринбург', 'https://images.unsplash.com/photo-1560250097-0b93528c311a?auto=format&fit=crop&w=300&q=80', 'Администратор площадки'),
  ('Алексей Власов', 'user@equipment.local', '$2a$10$zjVBxibw8De5HqJG5ftTr.lBT1sK/MoRJWr5h0oq99PA1QRIPqGvu', 'user', '+7 900 000-00-02', 'Екатеринбург', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=300&q=80', 'Сдаю технику для учебных и коммерческих съемок.')
ON CONFLICT (email) DO UPDATE SET role = EXCLUDED.role, password_hash = EXCLUDED.password_hash, blocked = false;

INSERT INTO categories (name, description)
VALUES
  ('Фото и видео', 'Камеры, объективы и стабилизаторы'),
  ('Презентации', 'Проекторы и оборудование для мероприятий'),
  ('Компьютеры', 'Ноутбуки и периферия')
ON CONFLICT (name) DO NOTHING;

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'user@equipment.local'), (SELECT id FROM categories WHERE name = 'Фото и видео'), 'Sony A7 III с объективом 24-70', 'Комплект для съемки: камера, объективы, аккумуляторы и сумка.', 'CAM-1001', 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32?auto=format&fit=crop&w=1200&q=80', 'Екатеринбург', 1800, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'CAM-1001');

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'user@equipment.local'), (SELECT id FROM categories WHERE name = 'Презентации'), 'Проектор Epson Full HD', 'Проектор для лекций, презентаций и домашних кинопоказов.', 'PRJ-2001', 'https://images.unsplash.com/photo-1601944179066-29786cb9d32a?auto=format&fit=crop&w=1200&q=80', 'Екатеринбург', 950, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'PRJ-2001');

INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available)
SELECT (SELECT id FROM users WHERE email = 'user@equipment.local'), (SELECT id FROM categories WHERE name = 'Фото и видео'), 'DJI Ronin RS 3', 'Стабилизатор для плавной видеосъемки и учебных проектов.', 'STB-3001', 'https://images.unsplash.com/photo-1616243850903-7b51366c1ba2?auto=format&fit=crop&w=1200&q=80', 'Пермь', 1200, true
WHERE NOT EXISTS (SELECT 1 FROM equipment WHERE serial = 'STB-3001');
