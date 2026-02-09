ALTER TABLE public.users
DROP COLUMN IF EXISTS username;

DROP INDEX IF EXISTS idx_users_username;
