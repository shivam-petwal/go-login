-- currencies table
CREATE TABLE IF NOT EXISTS public.currencies (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_currencies_code ON public.currencies (code);

--  exchange_rates table
CREATE TABLE IF NOT EXISTS public.exchange_rates (
    id BIGSERIAL PRIMARY KEY,
    from_currency_id BIGINT NOT NULL,
    to_currency_id BIGINT NOT NULL,
    rate DECIMAL(18, 6) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_from_currency FOREIGN KEY (from_currency_id) REFERENCES public.currencies (id) ON DELETE CASCADE,
    CONSTRAINT fk_to_currency FOREIGN KEY (to_currency_id) REFERENCES public.currencies (id) ON DELETE CASCADE
);


CREATE UNIQUE INDEX idx_exchange_rates_pair ON public.exchange_rates (from_currency_id, to_currency_id);