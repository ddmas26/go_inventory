ALTER TABLE users 
ADD COLUMN if not exists created_at TIMESTAMPTZ DEFAULT NOW(),
ADD COLUMN if not exists updated_at TIMESTAMPTZ;