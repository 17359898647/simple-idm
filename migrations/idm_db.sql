--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.4 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: idm
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


ALTER TABLE public.goose_db_version OWNER TO idm;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: idm
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.goose_db_version_id_seq OWNER TO idm;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: idm
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: idm
--

CREATE TABLE public.roles (
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    role_name character varying(255) NOT NULL,
    description text
);


ALTER TABLE public.roles OWNER TO idm;

--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: idm
--

CREATE TABLE public.user_roles (
    user_uuid uuid NOT NULL,
    role_uuid uuid NOT NULL
);


ALTER TABLE public.user_roles OWNER TO idm;

--
-- Name: users; Type: TABLE; Schema: public; Owner: idm
--

CREATE TABLE public.users (
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp without time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL,
    last_modified_at timestamp without time zone DEFAULT (now() AT TIME ZONE 'utc'::text) NOT NULL,
    deleted_at timestamp without time zone,
    created_by character varying(255),
    email character varying(255) NOT NULL,
    name character varying(255),
    password character varying(255),
    verified_at timestamp without time zone,
    username character varying(255)
);


ALTER TABLE public.users OWNER TO idm;

--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: idm
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: idm
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: idm
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (uuid);


--
-- Name: roles roles_role_name_key; Type: CONSTRAINT; Schema: public; Owner: idm
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_role_name_key UNIQUE (role_name);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: idm
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (user_uuid, role_uuid);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: idm
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (uuid);


--
-- Name: user_roles user_roles_role_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: idm
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_uuid_fkey FOREIGN KEY (role_uuid) REFERENCES public.roles(uuid);


--
-- Name: user_roles user_roles_user_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: idm
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_uuid_fkey FOREIGN KEY (user_uuid) REFERENCES public.users(uuid);


--
-- PostgreSQL database dump complete
--

