CREATE TABLE IF NOT EXISTS public.User (
    email VARCHAR(50) PRIMARY KEY,
    password VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS public.Item (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(50),
    index INT NOT NULL,
    count INT NOT NULL,
    FOREIGN KEY(user_email) REFERENCES public.User(email)
);


CREATE OR REPLACE PROCEDURE public.purchase_item (
    p_user_email VARCHAR(50),
    p_index INT,
    p_amount INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.Item 
        SET count = count + p_amount
        WHERE user_email = p_user_email AND index = p_index;

    IF NOT FOUND THEN
        INSERT INTO public.Item Values (DEFAULT, p_user_email, p_index, p_amount);
    END IF;
END;$$;