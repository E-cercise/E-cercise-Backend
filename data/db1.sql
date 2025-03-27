--
-- PostgreSQL database dump
--

-- Dumped from database version 16.8
-- Dumped by pg_dump version 16.8 (Homebrew)

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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: gender_type; Type: TYPE; Schema: public; Owner: -
--



--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: equipment; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: equipment_features; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: equipment_muscle_groups; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: equipment_options; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: goals; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: images; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: line_equipments; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: muscle_groups; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: -
--

--
-- Data for Name: user_preferences; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Name: attributes attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_pkey PRIMARY KEY (id);


--
-- Name: carts carts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_pkey PRIMARY KEY (id);


--
-- Name: equipment_features equipment_features_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.equipment_features
    ADD CONSTRAINT equipment_features_pkey PRIMARY KEY (id);


--
-- Name: equipment_muscle_groups equipment_muscle_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.equipment_muscle_groups
    ADD CONSTRAINT equipment_muscle_groups_pkey PRIMARY KEY (muscle_group_id, equipment_id);


--
-- Name: equipment_options equipment_options_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.equipment_options
    ADD CONSTRAINT equipment_options_pkey PRIMARY KEY (id);


--
-- Name: equipment equipment_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT equipment_pkey PRIMARY KEY (id);


--
-- Name: goals goals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goals
    ADD CONSTRAINT goals_pkey PRIMARY KEY (id);


--
-- Name: images images_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);


--
-- Name: line_equipments line_equipments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.line_equipments
    ADD CONSTRAINT line_equipments_pkey PRIMARY KEY (id);


--
-- Name: muscle_groups muscle_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.muscle_groups
    ADD CONSTRAINT muscle_groups_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: carts uni_carts_user_id; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT uni_carts_user_id UNIQUE (user_id);


--
-- Name: goals uni_goals_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goals
    ADD CONSTRAINT uni_goals_name UNIQUE (name);


--
-- Name: tags uni_tags_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT uni_tags_name UNIQUE (name);


--
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- Name: user_preferences user_preferences_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_preferences
    ADD CONSTRAINT user_preferences_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: line_equipments fk_carts_line_equipments; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.line_equipments
    ADD CONSTRAINT fk_carts_line_equipments FOREIGN KEY (cart_id) REFERENCES public.carts(id);


--
-- Name: attributes fk_equipment_attribute; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT fk_equipment_attribute FOREIGN KEY (equipment_id) REFERENCES public.equipment(id);


--
-- Name: equipment_features fk_equipment_equipment_feature; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.equipment_features
    ADD CONSTRAINT fk_equipment_equipment_feature FOREIGN KEY (equipment_id) REFERENCES public.equipment(id);


--
-- Name: equipment_options fk_equipment_equipment_options; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.equipment_options
    ADD CONSTRAINT fk_equipment_equipment_options FOREIGN KEY (equipment_id) REFERENCES public.equipment(id);


--
-- Name: equipment_muscle_groups fk_equipment_muscle_groups_equipment; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.equipment_muscle_groups
    ADD CONSTRAINT fk_equipment_muscle_groups_equipment FOREIGN KEY (equipment_id) REFERENCES public.equipment(id);


--
-- Name: equipment_muscle_groups fk_equipment_muscle_groups_muscle_group; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.equipment_muscle_groups
    ADD CONSTRAINT fk_equipment_muscle_groups_muscle_group FOREIGN KEY (muscle_group_id) REFERENCES public.muscle_groups(id);


--
-- Name: images fk_equipment_options_images; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.images
    ADD CONSTRAINT fk_equipment_options_images FOREIGN KEY (equipment_option_id) REFERENCES public.equipment_options(id);


--
-- Name: line_equipments fk_line_equipments_equipment_option; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.line_equipments
    ADD CONSTRAINT fk_line_equipments_equipment_option FOREIGN KEY (equipment_option_id) REFERENCES public.equipment_options(id);


--
-- Name: line_equipments fk_orders_line_equipments; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.line_equipments
    ADD CONSTRAINT fk_orders_line_equipments FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: user_preferences fk_user_preferences_tag; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_preferences
    ADD CONSTRAINT fk_user_preferences_tag FOREIGN KEY (tag_id) REFERENCES public.tags(id);


--
-- Name: carts fk_users_cart; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT fk_users_cart FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users fk_users_goal; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_goal FOREIGN KEY (goal_id) REFERENCES public.goals(id);


--
-- Name: orders fk_users_orders; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_users_orders FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_preferences fk_users_user_preferences; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_preferences
    ADD CONSTRAINT fk_users_user_preferences FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

