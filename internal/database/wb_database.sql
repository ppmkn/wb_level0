--
-- PostgreSQL database dump
--

-- Dumped from database version 16.0
-- Dumped by pg_dump version 16.0

-- Started on 2023-10-20 22:16:13

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
-- TOC entry 216 (class 1259 OID 16510)
-- Name: delivery; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.delivery (
    order_uid text,
    name text,
    phone text,
    zip text,
    city text,
    address text,
    region text,
    email text
);


ALTER TABLE public.delivery OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 16530)
-- Name: items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.items (
    order_uid text,
    chrt_id integer NOT NULL,
    track_number text,
    price numeric,
    rid text,
    name text,
    sale numeric,
    size text,
    total_price numeric,
    nm_id integer,
    brand text,
    status integer
);


ALTER TABLE public.items OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 16503)
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    order_uid text NOT NULL,
    track_number text,
    entry text,
    locale text,
    internal_signature text,
    customer_id text,
    delivery_service text,
    shardkey text,
    sm_id integer,
    oof_shard text,
    date_created timestamp without time zone
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16520)
-- Name: payment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment (
    order_uid text,
    transaction text,
    request_id text,
    currency text,
    provider text,
    amount numeric,
    payment_dt timestamp without time zone,
    bank text,
    delivery_cost numeric,
    goods_total numeric,
    custom_fee numeric
);


ALTER TABLE public.payment OWNER TO postgres;

--
-- TOC entry 4796 (class 0 OID 16510)
-- Dependencies: 216
-- Data for Name: delivery; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.delivery (order_uid, name, phone, zip, city, address, region, email) FROM stdin;
b563feb7b2b84b6test	Test Testov	+9720000000	2639809	Kiryat Mozkin	Ploshad Mira 15	Kraiot	test@gmail.com
\.


--
-- TOC entry 4798 (class 0 OID 16530)
-- Dependencies: 218
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) FROM stdin;
b563feb7b2b84b6test	9934930	WBILMTESTTRACK	453	ab4219087a764ae0btest	Mascaras	30	0	317	2389212	Vivienne Sabo	202
\.


--
-- TOC entry 4795 (class 0 OID 16503)
-- Dependencies: 215
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, oof_shard, date_created) FROM stdin;
b563feb7b2b84b6test	WBILMTESTTRACK	WBIL	en		test	meest	9	99	1	2021-11-26 06:22:19
b222feb7b2b84b6test	WBILMTESTTRACK	WBIL	en		test	meest	9	99	1	2021-11-26 06:22:19
a123feb7b2b84b6test	WBILMTESTTRACK	WBIL	\N	\N	\N	\N	\N	\N	\N	\N
a234feb7b2b84b6test	WBILMTESTTRACK	WBIL	en		test	meest	9	99	1	2021-11-26 06:22:19
y787feb7b2b84b6test	WBILMTESTTRACK	WBIL	en		test	meest	9	99	1	2021-11-26 06:22:19
z222feb7b2b84b6test	WBILMTESTTRACK	WBIL	\N	\N	\N	\N	\N	\N	\N	\N
23g4234234234234234	23g423423423423423423g423423423423423423g4223g423423423423423423g423423423423423434234234234234	№;;№%	\N	\N	\N	\N	\N	\N	\N	\N
erdfgdfgdfg	WBILMTESTTRACK	WBIL	\N	\N	\N	\N	\N	\N	\N	\N
\.


--
-- TOC entry 4797 (class 0 OID 16520)
-- Dependencies: 217
-- Data for Name: payment; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) FROM stdin;
b563feb7b2b84b6test	b563feb7b2b84b6test		USD	wbpay	1817	2023-10-10 12:00:00	alpha	1500	317	0
\.


--
-- TOC entry 4648 (class 2606 OID 16536)
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (chrt_id);


--
-- TOC entry 4646 (class 2606 OID 16509)
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (order_uid);


--
-- TOC entry 4649 (class 2606 OID 16515)
-- Name: delivery delivery_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.delivery
    ADD CONSTRAINT delivery_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid);


--
-- TOC entry 4651 (class 2606 OID 16537)
-- Name: items items_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid);


--
-- TOC entry 4650 (class 2606 OID 16525)
-- Name: payment payment_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid);


-- Completed on 2023-10-20 22:16:13

--
-- PostgreSQL database dump complete
--

