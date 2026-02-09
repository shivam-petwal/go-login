-- Add deleted_at column to currencies table
ALTER TABLE public.currencies ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX idx_currencies_deleted_at ON public.currencies (deleted_at);

-- -- Add deleted_at column to exchange_rates table
ALTER TABLE public.exchange_rates ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX idx_exchange_rates_deleted_at ON public.exchange_rates (deleted_at);
