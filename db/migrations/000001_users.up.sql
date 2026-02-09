
ALTER TABLE public.users
ADD COLUMN username VARCHAR(50);


UPDATE public.users
SET username = CONCAT('user_', id)
WHERE username IS NULL;


ALTER TABLE public.users
ALTER COLUMN username SET NOT NULL;


CREATE UNIQUE INDEX idx_users_username
ON public.users (username);
