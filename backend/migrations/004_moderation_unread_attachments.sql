DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'moderator' AND enumtypid = 'user_role'::regtype) THEN
    ALTER TYPE user_role ADD VALUE 'moderator';
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_status') THEN
    CREATE TYPE moderation_status AS ENUM ('pending', 'approved', 'rejected');
  END IF;
END $$;

ALTER TABLE equipment ADD COLUMN IF NOT EXISTS moderation_status moderation_status NOT NULL DEFAULT 'approved';
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS reject_reason TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN IF NOT EXISTS attachment_url TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN IF NOT EXISTS read_at TIMESTAMPTZ;

INSERT INTO users (name, email, password_hash, role, phone, city, avatar_url, bio)
VALUES
  ('Модератор площадки', 'moderator@equipment.local', '$2a$10$zjVBxibw8De5HqJG5ftTr.lBT1sK/MoRJWr5h0oq99PA1QRIPqGvu', 'moderator', '+7 900 000-00-06', 'Екатеринбург', 'https://images.unsplash.com/photo-1580894908361-967195033215?auto=format&fit=crop&w=300&q=80', 'Проверяю объявления перед публикацией.')
ON CONFLICT (email) DO UPDATE SET role = EXCLUDED.role, password_hash = EXCLUDED.password_hash, blocked = false;

UPDATE equipment SET moderation_status = 'approved' WHERE moderation_status IS NULL;
