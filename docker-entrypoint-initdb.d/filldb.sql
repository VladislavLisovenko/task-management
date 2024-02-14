CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    name character varying(50) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.tasks
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    description character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    user_id integer NOT NULL,
    expiration_date date NOT NULL DEFAULT now(),
    done boolean,
    CONSTRAINT task_pkey PRIMARY KEY (id),
    CONSTRAINT tasks_users FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tasks
    OWNER to postgres;

delete from public.tasks;
delete from public.users;

ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE tasks_id_seq RESTART WITH 1;

INSERT INTO public.users(name)
VALUES 
    ('Dmitriy'),
	('Andrey'),
	('Nikolay');
	
INSERT INTO public.tasks(description, user_id, expiration_date, done)
VALUES 
    ('To wash a car', 1, '2024/01/31', false),
    ('To paint a fence', 2, '2024/02/28', false),
    ('To start study playing a piano', 3, '2024/03/31', false);