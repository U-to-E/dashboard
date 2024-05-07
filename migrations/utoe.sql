--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2
-- Dumped by pg_dump version 16.2

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
-- Name: U100; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."U100" (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    email text,
    collage_id text NOT NULL,
    mentor_id text NOT NULL,
    level bigint
);


ALTER TABLE public."U100" OWNER TO postgres;

--
-- Name: U100_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."U100_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."U100_id_seq" OWNER TO postgres;

--
-- Name: U100_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."U100_id_seq" OWNED BY public."U100".id;


--
-- Name: logins; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.logins (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    email text NOT NULL,
    password bytea NOT NULL,
    collage_id text,
    role text
);


ALTER TABLE public.logins OWNER TO postgres;

--
-- Name: logins_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.logins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.logins_id_seq OWNER TO postgres;

--
-- Name: logins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.logins_id_seq OWNED BY public.logins.id;


--
-- Name: mappings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mappings (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    mentor_id text NOT NULL,
    college_id text NOT NULL,
    students_id integer[]
);


ALTER TABLE public.mappings OWNER TO postgres;

--
-- Name: mappings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.mappings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mappings_id_seq OWNER TO postgres;

--
-- Name: mappings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.mappings_id_seq OWNED BY public.mappings.id;


--
-- Name: mentors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mentors (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    email text NOT NULL
);


ALTER TABLE public.mentors OWNER TO postgres;

--
-- Name: mentors_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.mentors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mentors_id_seq OWNER TO postgres;

--
-- Name: mentors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.mentors_id_seq OWNED BY public.mentors.id;


--
-- Name: students; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.students (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    email text,
    collage_id text NOT NULL,
    mentor_id text NOT NULL,
    level bigint
);


ALTER TABLE public.students OWNER TO postgres;

--
-- Name: students_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.students_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.students_id_seq OWNER TO postgres;

--
-- Name: students_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.students_id_seq OWNED BY public.students.id;


--
-- Name: U100 id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."U100" ALTER COLUMN id SET DEFAULT nextval('public."U100_id_seq"'::regclass);


--
-- Name: logins id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logins ALTER COLUMN id SET DEFAULT nextval('public.logins_id_seq'::regclass);


--
-- Name: mappings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mappings ALTER COLUMN id SET DEFAULT nextval('public.mappings_id_seq'::regclass);


--
-- Name: mentors id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mentors ALTER COLUMN id SET DEFAULT nextval('public.mentors_id_seq'::regclass);


--
-- Name: students id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students ALTER COLUMN id SET DEFAULT nextval('public.students_id_seq'::regclass);


--
-- Data for Name: U100; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."U100" (id, created_at, updated_at, deleted_at, name, email, collage_id, mentor_id, level) FROM stdin;
1	2024-05-04 16:16:33.391185+05:30	2024-05-04 16:16:33.391185+05:30	\N	Harshith reddy	hmadhadi@gitam.in	U100	0	1
2	2024-05-04 16:16:33.4123+05:30	2024-05-04 16:16:33.4123+05:30	\N	Rahul	hsdbfh@hh.in	U100	0	1
3	2024-05-04 16:19:06.551096+05:30	2024-05-04 16:19:06.551096+05:30	\N	Harshith reddy	hmadhadi@gitam.in	U100	0	1
4	2024-05-04 16:19:06.55645+05:30	2024-05-04 16:19:06.55645+05:30	\N	Rahul	hsdbfh@hh.in	U100	0	1
5	2024-05-04 16:19:06.561289+05:30	2024-05-04 16:19:06.561289+05:30	\N	Big	jkdfh@jjdb.jj	U100	0	1
6	2024-05-04 21:48:10.001381+05:30	2024-05-04 21:48:10.001381+05:30	\N	Harshith reddy	hmadhadi@gitam.in	U100	0	1
7	2024-05-04 21:48:10.763694+05:30	2024-05-04 21:48:10.763694+05:30	\N	Rahul	hsdbfh@hh.in	U100	0	1
8	2024-05-04 21:48:11.512087+05:30	2024-05-04 21:48:11.512087+05:30	\N	Big	jkdfh@jjdb.jj	U100	0	1
9	2024-05-04 21:52:38.212128+05:30	2024-05-04 21:52:38.212128+05:30	\N	Harshith reddy	hmadhadi@gitam.in	U100	0	1
10	2024-05-04 21:52:38.946268+05:30	2024-05-04 21:52:38.946268+05:30	\N	Rahul	hsdbfh@hh.in	U100	0	1
11	2024-05-04 21:52:39.677767+05:30	2024-05-04 21:52:39.677767+05:30	\N	Big	jkdfh@jjdb.jj	U100	0	1
12	2024-05-04 21:54:50.830606+05:30	2024-05-04 21:54:50.830606+05:30	\N	Harshith reddy	hmadhadi@gitam.in	U100	0	1
13	2024-05-04 21:54:51.55928+05:30	2024-05-04 21:54:51.55928+05:30	\N	Rahul	hsdbfh@hh.in	U100	0	1
14	2024-05-04 21:54:52.288784+05:30	2024-05-04 21:54:52.288784+05:30	\N	Big	jkdfh@jjdb.jj	U100	0	1
15	2024-05-04 21:56:10.986779+05:30	2024-05-04 21:56:10.986779+05:30	\N	Harshith reddy	hmadhadi@gitam.in	U100	0	1
16	2024-05-04 21:56:11.72926+05:30	2024-05-04 21:56:11.72926+05:30	\N	Rahul	hsdbfh@hh.in	U100	0	1
17	2024-05-04 21:56:12.472366+05:30	2024-05-04 21:56:12.472366+05:30	\N	Big	jkdfh@jjdb.jj	U100	0	1
18	2024-05-04 22:07:08.943078+05:30	2024-05-04 22:07:08.943078+05:30	\N	Harshith reddy	hmadhadi@gitam.in	U100	0	1
19	2024-05-04 22:07:08.94819+05:30	2024-05-04 22:07:08.94819+05:30	\N	Rahul	hsdbfh@hh.in	U100	0	1
20	2024-05-04 22:07:08.953489+05:30	2024-05-04 22:07:08.953489+05:30	\N	Big	jkdfh@jjdb.jj	U100	0	1
\.


--
-- Data for Name: logins; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.logins (id, created_at, updated_at, deleted_at, name, email, password, collage_id, role) FROM stdin;
1	2024-05-04 22:07:08.938905+05:30	2024-05-04 22:07:08.938905+05:30	\N	Harshith reddy	hmadhadi@gitam.in	\\x2432612431342461633343766b4f552f4c6172543645365066562e64656d476c7376304c72515642704d364f546d49766b684e585373377363744536	U100	Student
2	2024-05-04 22:07:08.945427+05:30	2024-05-04 22:07:08.945427+05:30	\N	Rahul	hsdbfh@hh.in	\\x2432612431342461633343766b4f552f4c6172543645365066562e64656d476c7376304c72515642704d364f546d49766b684e585373377363744536	U100	Student
3	2024-05-04 22:07:08.950498+05:30	2024-05-04 22:07:08.950498+05:30	\N	Big	jkdfh@jjdb.jj	\\x2432612431342461633343766b4f552f4c6172543645365066562e64656d476c7376304c72515642704d364f546d49766b684e585373377363744536	U100	Student
4	2024-05-06 18:15:43.093796+05:30	2024-05-06 18:15:43.093796+05:30	\N	Harish	hrish@utoe.in	\\x243261243134243047526b387742396d4d6b77574d784e474b644f4e2e2f596d7a596744653273796a4d7a42374a6b494f35422f3073433277625875		Mentor
5	2024-05-06 18:16:11.134992+05:30	2024-05-06 18:16:11.134992+05:30	\N	Home Server	hiii@adm.in	\\x243261243134246b477254774c58554153432e48664e4f67616334466542495141474a49737868305766795a5a384e36316e656a4e767a7334736936	U200	Student
7	2024-05-06 18:18:54.729662+05:30	2024-05-06 18:18:54.729662+05:30	\N	Home Server	hiii@adm.com	\\x243261243134246b477254774c58554153432e48664e4f67616334466542495141474a49737868305766795a5a384e36316e656a4e767a7334736936	U200	Student
\.


--
-- Data for Name: mappings; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mappings (id, created_at, updated_at, deleted_at, mentor_id, college_id, students_id) FROM stdin;
\.


--
-- Data for Name: mentors; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mentors (id, created_at, updated_at, deleted_at, name, email) FROM stdin;
1	2024-05-06 18:15:43.100093+05:30	2024-05-06 18:15:43.100093+05:30	\N	Harish	hrish@utoe.in
\.


--
-- Data for Name: students; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.students (id, created_at, updated_at, deleted_at, name, email, collage_id, mentor_id, level) FROM stdin;
\.


--
-- Name: U100_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."U100_id_seq"', 20, true);


--
-- Name: logins_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.logins_id_seq', 8, true);


--
-- Name: mappings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.mappings_id_seq', 1, false);


--
-- Name: mentors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.mentors_id_seq', 1, true);


--
-- Name: students_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.students_id_seq', 1, false);


--
-- Name: U100 U100_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."U100"
    ADD CONSTRAINT "U100_pkey" PRIMARY KEY (id);


--
-- Name: logins logins_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logins
    ADD CONSTRAINT logins_pkey PRIMARY KEY (id);


--
-- Name: mappings mappings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mappings
    ADD CONSTRAINT mappings_pkey PRIMARY KEY (id);


--
-- Name: mentors mentors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mentors
    ADD CONSTRAINT mentors_pkey PRIMARY KEY (id);


--
-- Name: students students_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_pkey PRIMARY KEY (id);


--
-- Name: idx_U100_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "idx_U100_deleted_at" ON public."U100" USING btree (deleted_at);


--
-- Name: idx_logins_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_logins_deleted_at ON public.logins USING btree (deleted_at);


--
-- Name: idx_logins_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_logins_email ON public.logins USING btree (email);


--
-- Name: idx_mappings_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_mappings_deleted_at ON public.mappings USING btree (deleted_at);


--
-- Name: idx_mentors_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_mentors_deleted_at ON public.mentors USING btree (deleted_at);


--
-- Name: idx_students_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_students_deleted_at ON public.students USING btree (deleted_at);


--
-- PostgreSQL database dump complete
--

