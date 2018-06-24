--
-- PostgreSQL database dump
--

-- Dumped from database version 10.4
-- Dumped by pg_dump version 10.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: stream; Type: TABLE; Schema: public; Owner: testu
--

CREATE TABLE public.stream (
    id character varying(255) NOT NULL,
    status character varying(255) NOT NULL
);


ALTER TABLE public.stream OWNER TO testu;

--
-- Data for Name: stream; Type: TABLE DATA; Schema: public; Owner: testu
--

COPY public.stream (id, status) FROM stdin;
07173261-97ee-412a-b310-652af8c3a52e	Created
85379596-2768-4974-957d-c8cda780300a	Created
14d5a503-41cb-4b89-bf50-5431f7be093e	Created
\.

--
-- PostgreSQL database dump complete
--

