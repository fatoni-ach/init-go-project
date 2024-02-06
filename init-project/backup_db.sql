COMMENT ON SCHEMA public IS 'standard public schema';

-- DROP SEQUENCE public.user_id_seq;

CREATE SEQUENCE public.user_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id int4 NOT NULL DEFAULT nextval('user_id_seq'::regclass),
	username varchar NOT NULL,
	email varchar NOT NULL,
	"password" varchar NOT NULL,
	fullname varchar NOT NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	gender bool NULL,
	deleted_at timestamp NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);