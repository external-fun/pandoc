--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2 (Debian 15.2-1.pgdg110+1)
-- Dumped by pg_dump version 15.2

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

ALTER DATABASE pandoc_db OWNER TO lispberry;

\connect pandoc_db

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: file; Type: TABLE; Schema: public; Owner: lispberry
--

CREATE TABLE public.file (
                             uuid character varying NOT NULL,
                             origin_url character varying,
                             converted_url character varying,
                             status character varying
);


ALTER TABLE public.file OWNER TO lispberry;

--
-- Data for Name: file; Type: TABLE DATA; Schema: public; Owner: lispberry
--

COPY public.file (uuid, origin_url, converted_url, status) FROM stdin;
\.


--
-- Name: file file_converted_url_key; Type: CONSTRAINT; Schema: public; Owner: lispberry
--

ALTER TABLE ONLY public.file
    ADD CONSTRAINT file_converted_url_key UNIQUE (converted_url);


--
-- Name: file file_origin_url_key; Type: CONSTRAINT; Schema: public; Owner: lispberry
--

ALTER TABLE ONLY public.file
    ADD CONSTRAINT file_origin_url_key UNIQUE (origin_url);


--
-- Name: file file_pkey; Type: CONSTRAINT; Schema: public; Owner: lispberry
--

ALTER TABLE ONLY public.file
    ADD CONSTRAINT file_pkey PRIMARY KEY (uuid);


--
-- PostgreSQL database dump complete
--

