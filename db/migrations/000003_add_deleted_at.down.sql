
ALTER TABLE public.exchange_rates DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE public.currencies DROP COLUMN IF EXISTS deleted_at;
