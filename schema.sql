--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6
-- Dumped by pg_dump version 17.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: _locales; Type: TYPE; Schema: public; Owner: mensauser
--

CREATE TYPE public._locales AS ENUM (
    'de',
    'en'
);


ALTER TYPE public._locales OWNER TO mensauser;

--
-- Name: enum_locale_locale; Type: TYPE; Schema: public; Owner: mensauser
--

CREATE TYPE public.enum_locale_locale AS ENUM (
    'de',
    'en'
);


ALTER TYPE public.enum_locale_locale OWNER TO mensauser;

--
-- Name: enum_recipes_category; Type: TYPE; Schema: public; Owner: mensauser
--

CREATE TYPE public.enum_recipes_category AS ENUM (
    'starter',
    'main',
    'side',
    'dessert'
);


ALTER TYPE public.enum_recipes_category OWNER TO mensauser;

--
-- Name: enum_serving_time_day; Type: TYPE; Schema: public; Owner: mensauser
--

CREATE TYPE public.enum_serving_time_day AS ENUM (
    'monday',
    'tuesday',
    'wednesday',
    'thursday',
    'friday',
    'saturday',
    'sunday'
);


ALTER TYPE public.enum_serving_time_day OWNER TO mensauser;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: additives; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.additives (
    id integer NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.additives OWNER TO mensauser;

--
-- Name: additives_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.additives_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.additives_id_seq OWNER TO mensauser;

--
-- Name: additives_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.additives_id_seq OWNED BY public.additives.id;


--
-- Name: additives_locales; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.additives_locales (
    name character varying,
    id integer NOT NULL,
    _locale public._locales NOT NULL,
    _parent_id integer NOT NULL
);


ALTER TABLE public.additives_locales OWNER TO mensauser;

--
-- Name: additives_locales_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.additives_locales_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.additives_locales_id_seq OWNER TO mensauser;

--
-- Name: additives_locales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.additives_locales_id_seq OWNED BY public.additives_locales.id;


--
-- Name: allergens; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.allergens (
    id integer NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.allergens OWNER TO mensauser;

--
-- Name: allergens_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.allergens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.allergens_id_seq OWNER TO mensauser;

--
-- Name: allergens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.allergens_id_seq OWNED BY public.allergens.id;


--
-- Name: allergens_locales; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.allergens_locales (
    name character varying,
    id integer NOT NULL,
    _locale public._locales NOT NULL,
    _parent_id integer NOT NULL
);


ALTER TABLE public.allergens_locales OWNER TO mensauser;

--
-- Name: allergens_locales_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.allergens_locales_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.allergens_locales_id_seq OWNER TO mensauser;

--
-- Name: allergens_locales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.allergens_locales_id_seq OWNED BY public.allergens_locales.id;


--
-- Name: features; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.features (
    id integer NOT NULL,
    visible_small boolean DEFAULT false,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.features OWNER TO mensauser;

--
-- Name: features_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.features_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.features_id_seq OWNER TO mensauser;

--
-- Name: features_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.features_id_seq OWNED BY public.features.id;


--
-- Name: info; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.info (
    id integer NOT NULL,
    title character varying,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.info OWNER TO mensauser;

--
-- Name: info_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.info_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.info_id_seq OWNER TO mensauser;

--
-- Name: info_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.info_id_seq OWNED BY public.info.id;


--
-- Name: locale; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.locale (
    id integer NOT NULL,
    name character varying NOT NULL,
    locale public.enum_locale_locale NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.locale OWNER TO mensauser;

--
-- Name: locale_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.locale_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.locale_id_seq OWNER TO mensauser;

--
-- Name: locale_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.locale_id_seq OWNED BY public.locale.id;


--
-- Name: locale_rels; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.locale_rels (
    id integer NOT NULL,
    "order" integer,
    parent_id integer NOT NULL,
    path character varying NOT NULL,
    recipes_id integer,
    features_id integer,
    allergens_id integer,
    additives_id integer
);


ALTER TABLE public.locale_rels OWNER TO mensauser;

--
-- Name: locale_rels_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.locale_rels_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.locale_rels_id_seq OWNER TO mensauser;

--
-- Name: locale_rels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.locale_rels_id_seq OWNED BY public.locale_rels.id;


--
-- Name: media; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.media (
    id integer NOT NULL,
    alt character varying NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    url character varying,
    thumbnail_u_r_l character varying,
    filename character varying,
    mime_type character varying,
    filesize numeric,
    width numeric,
    height numeric,
    focal_x numeric,
    focal_y numeric
);


ALTER TABLE public.media OWNER TO mensauser;

--
-- Name: media_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.media_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.media_id_seq OWNER TO mensauser;

--
-- Name: media_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.media_id_seq OWNED BY public.media.id;


--
-- Name: mensa; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.mensa (
    id integer NOT NULL,
    name character varying NOT NULL,
    slug character varying,
    address_latitude numeric NOT NULL,
    address_longitude numeric NOT NULL,
    address_street character varying,
    address_house_number character varying,
    address_zip_code character varying,
    address_city character varying,
    description character varying,
    provider_id integer NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.mensa OWNER TO mensauser;

--
-- Name: mensa_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.mensa_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mensa_id_seq OWNER TO mensauser;

--
-- Name: mensa_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.mensa_id_seq OWNED BY public.mensa.id;


--
-- Name: mensa_provider; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.mensa_provider (
    id integer NOT NULL,
    name character varying NOT NULL,
    slug character varying,
    description character varying NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.mensa_provider OWNER TO mensauser;

--
-- Name: mensa_provider_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.mensa_provider_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mensa_provider_id_seq OWNER TO mensauser;

--
-- Name: mensa_provider_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.mensa_provider_id_seq OWNED BY public.mensa_provider.id;


--
-- Name: nutrient_labels; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.nutrient_labels (
    id integer NOT NULL,
    unit_id integer NOT NULL,
    recommendation character varying,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.nutrient_labels OWNER TO mensauser;

--
-- Name: nutrient_labels_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.nutrient_labels_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.nutrient_labels_id_seq OWNER TO mensauser;

--
-- Name: nutrient_labels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.nutrient_labels_id_seq OWNED BY public.nutrient_labels.id;


--
-- Name: nutrient_labels_locales; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.nutrient_labels_locales (
    name character varying NOT NULL,
    id integer NOT NULL,
    _locale public._locales NOT NULL,
    _parent_id integer NOT NULL
);


ALTER TABLE public.nutrient_labels_locales OWNER TO mensauser;

--
-- Name: nutrient_labels_locales_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.nutrient_labels_locales_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.nutrient_labels_locales_id_seq OWNER TO mensauser;

--
-- Name: nutrient_labels_locales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.nutrient_labels_locales_id_seq OWNED BY public.nutrient_labels_locales.id;


--
-- Name: nutrient_units; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.nutrient_units (
    id integer NOT NULL,
    name character varying NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.nutrient_units OWNER TO mensauser;

--
-- Name: nutrient_units_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.nutrient_units_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.nutrient_units_id_seq OWNER TO mensauser;

--
-- Name: nutrient_units_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.nutrient_units_id_seq OWNED BY public.nutrient_units.id;


--
-- Name: nutrient_values; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.nutrient_values (
    id integer NOT NULL,
    value numeric,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.nutrient_values OWNER TO mensauser;

--
-- Name: nutrient_values_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.nutrient_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.nutrient_values_id_seq OWNER TO mensauser;

--
-- Name: nutrient_values_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.nutrient_values_id_seq OWNED BY public.nutrient_values.id;


--
-- Name: nutrients; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.nutrients (
    id integer NOT NULL,
    nutrient_value_id integer,
    nutrient_label_id integer,
    recipe_id integer,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.nutrients OWNER TO mensauser;

--
-- Name: nutrients_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.nutrients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.nutrients_id_seq OWNER TO mensauser;

--
-- Name: nutrients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.nutrients_id_seq OWNED BY public.nutrients.id;


--
-- Name: payload_locked_documents; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.payload_locked_documents (
    id integer NOT NULL,
    global_slug character varying,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.payload_locked_documents OWNER TO mensauser;

--
-- Name: payload_locked_documents_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.payload_locked_documents_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payload_locked_documents_id_seq OWNER TO mensauser;

--
-- Name: payload_locked_documents_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.payload_locked_documents_id_seq OWNED BY public.payload_locked_documents.id;


--
-- Name: payload_locked_documents_rels; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.payload_locked_documents_rels (
    id integer NOT NULL,
    "order" integer,
    parent_id integer NOT NULL,
    path character varying NOT NULL,
    users_id integer,
    media_id integer,
    mensa_provider_id integer,
    mensa_id integer,
    serving_time_id integer,
    time_slot_id integer,
    servings_id integer,
    info_id integer,
    nutrients_id integer,
    nutrient_units_id integer,
    nutrient_labels_id integer,
    nutrient_values_id integer,
    allergens_id integer,
    additives_id integer,
    recipes_id integer,
    user_image_uploads_id integer,
    features_id integer,
    locale_id integer
);


ALTER TABLE public.payload_locked_documents_rels OWNER TO mensauser;

--
-- Name: payload_locked_documents_rels_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.payload_locked_documents_rels_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payload_locked_documents_rels_id_seq OWNER TO mensauser;

--
-- Name: payload_locked_documents_rels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.payload_locked_documents_rels_id_seq OWNED BY public.payload_locked_documents_rels.id;


--
-- Name: payload_migrations; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.payload_migrations (
    id integer NOT NULL,
    name character varying,
    batch numeric,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.payload_migrations OWNER TO mensauser;

--
-- Name: payload_migrations_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.payload_migrations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payload_migrations_id_seq OWNER TO mensauser;

--
-- Name: payload_migrations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.payload_migrations_id_seq OWNED BY public.payload_migrations.id;


--
-- Name: payload_preferences; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.payload_preferences (
    id integer NOT NULL,
    key character varying,
    value jsonb,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.payload_preferences OWNER TO mensauser;

--
-- Name: payload_preferences_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.payload_preferences_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payload_preferences_id_seq OWNER TO mensauser;

--
-- Name: payload_preferences_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.payload_preferences_id_seq OWNED BY public.payload_preferences.id;


--
-- Name: payload_preferences_rels; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.payload_preferences_rels (
    id integer NOT NULL,
    "order" integer,
    parent_id integer NOT NULL,
    path character varying NOT NULL,
    users_id integer
);


ALTER TABLE public.payload_preferences_rels OWNER TO mensauser;

--
-- Name: payload_preferences_rels_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.payload_preferences_rels_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payload_preferences_rels_id_seq OWNER TO mensauser;

--
-- Name: payload_preferences_rels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.payload_preferences_rels_id_seq OWNED BY public.payload_preferences_rels.id;


--
-- Name: recipes; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.recipes (
    id integer NOT NULL,
    ai_thumbnail_id integer,
    price_students numeric,
    price_employees numeric,
    price_guests numeric,
    mensa_provider_id integer NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    category public.enum_recipes_category
);


ALTER TABLE public.recipes OWNER TO mensauser;

--
-- Name: recipes_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.recipes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recipes_id_seq OWNER TO mensauser;

--
-- Name: recipes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.recipes_id_seq OWNED BY public.recipes.id;


--
-- Name: recipes_rels; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.recipes_rels (
    id integer NOT NULL,
    "order" integer,
    parent_id integer NOT NULL,
    path character varying NOT NULL,
    features_id integer,
    additives_id integer,
    allergens_id integer
);


ALTER TABLE public.recipes_rels OWNER TO mensauser;

--
-- Name: recipes_rels_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.recipes_rels_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recipes_rels_id_seq OWNER TO mensauser;

--
-- Name: recipes_rels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.recipes_rels_id_seq OWNED BY public.recipes_rels.id;


--
-- Name: serving_time; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.serving_time (
    id integer NOT NULL,
    mensa_info_id integer,
    day public.enum_serving_time_day,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.serving_time OWNER TO mensauser;

--
-- Name: serving_time_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.serving_time_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.serving_time_id_seq OWNER TO mensauser;

--
-- Name: serving_time_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.serving_time_id_seq OWNED BY public.serving_time.id;


--
-- Name: servings; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.servings (
    id integer NOT NULL,
    recipe_id integer NOT NULL,
    date timestamp(3) with time zone NOT NULL,
    mensa_id integer,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.servings OWNER TO mensauser;

--
-- Name: servings_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.servings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.servings_id_seq OWNER TO mensauser;

--
-- Name: servings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.servings_id_seq OWNED BY public.servings.id;


--
-- Name: time_slot; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.time_slot (
    id integer NOT NULL,
    serving_time_id integer,
    "from" timestamp(3) with time zone NOT NULL,
    "to" timestamp(3) with time zone NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.time_slot OWNER TO mensauser;

--
-- Name: time_slot_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.time_slot_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.time_slot_id_seq OWNER TO mensauser;

--
-- Name: time_slot_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.time_slot_id_seq OWNED BY public.time_slot.id;


--
-- Name: user_image_uploads; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.user_image_uploads (
    id integer NOT NULL,
    image_id integer NOT NULL,
    unique_user_id character varying NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    recipe_id integer NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.user_image_uploads OWNER TO mensauser;

--
-- Name: user_image_uploads_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.user_image_uploads_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_image_uploads_id_seq OWNER TO mensauser;

--
-- Name: user_image_uploads_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.user_image_uploads_id_seq OWNED BY public.user_image_uploads.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: mensauser
--

CREATE TABLE public.users (
    id integer NOT NULL,
    updated_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    email character varying NOT NULL,
    reset_password_token character varying,
    reset_password_expiration timestamp(3) with time zone,
    salt character varying,
    hash character varying,
    login_attempts numeric DEFAULT 0,
    lock_until timestamp(3) with time zone
);


ALTER TABLE public.users OWNER TO mensauser;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: mensauser
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO mensauser;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mensauser
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: additives id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.additives ALTER COLUMN id SET DEFAULT nextval('public.additives_id_seq'::regclass);


--
-- Name: additives_locales id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.additives_locales ALTER COLUMN id SET DEFAULT nextval('public.additives_locales_id_seq'::regclass);


--
-- Name: allergens id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.allergens ALTER COLUMN id SET DEFAULT nextval('public.allergens_id_seq'::regclass);


--
-- Name: allergens_locales id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.allergens_locales ALTER COLUMN id SET DEFAULT nextval('public.allergens_locales_id_seq'::regclass);


--
-- Name: features id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.features ALTER COLUMN id SET DEFAULT nextval('public.features_id_seq'::regclass);


--
-- Name: info id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.info ALTER COLUMN id SET DEFAULT nextval('public.info_id_seq'::regclass);


--
-- Name: locale id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale ALTER COLUMN id SET DEFAULT nextval('public.locale_id_seq'::regclass);


--
-- Name: locale_rels id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale_rels ALTER COLUMN id SET DEFAULT nextval('public.locale_rels_id_seq'::regclass);


--
-- Name: media id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.media ALTER COLUMN id SET DEFAULT nextval('public.media_id_seq'::regclass);


--
-- Name: mensa id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.mensa ALTER COLUMN id SET DEFAULT nextval('public.mensa_id_seq'::regclass);


--
-- Name: mensa_provider id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.mensa_provider ALTER COLUMN id SET DEFAULT nextval('public.mensa_provider_id_seq'::regclass);


--
-- Name: nutrient_labels id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_labels ALTER COLUMN id SET DEFAULT nextval('public.nutrient_labels_id_seq'::regclass);


--
-- Name: nutrient_labels_locales id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_labels_locales ALTER COLUMN id SET DEFAULT nextval('public.nutrient_labels_locales_id_seq'::regclass);


--
-- Name: nutrient_units id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_units ALTER COLUMN id SET DEFAULT nextval('public.nutrient_units_id_seq'::regclass);


--
-- Name: nutrient_values id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_values ALTER COLUMN id SET DEFAULT nextval('public.nutrient_values_id_seq'::regclass);


--
-- Name: nutrients id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrients ALTER COLUMN id SET DEFAULT nextval('public.nutrients_id_seq'::regclass);


--
-- Name: payload_locked_documents id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents ALTER COLUMN id SET DEFAULT nextval('public.payload_locked_documents_id_seq'::regclass);


--
-- Name: payload_locked_documents_rels id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels ALTER COLUMN id SET DEFAULT nextval('public.payload_locked_documents_rels_id_seq'::regclass);


--
-- Name: payload_migrations id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_migrations ALTER COLUMN id SET DEFAULT nextval('public.payload_migrations_id_seq'::regclass);


--
-- Name: payload_preferences id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_preferences ALTER COLUMN id SET DEFAULT nextval('public.payload_preferences_id_seq'::regclass);


--
-- Name: payload_preferences_rels id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_preferences_rels ALTER COLUMN id SET DEFAULT nextval('public.payload_preferences_rels_id_seq'::regclass);


--
-- Name: recipes id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes ALTER COLUMN id SET DEFAULT nextval('public.recipes_id_seq'::regclass);


--
-- Name: recipes_rels id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes_rels ALTER COLUMN id SET DEFAULT nextval('public.recipes_rels_id_seq'::regclass);


--
-- Name: serving_time id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.serving_time ALTER COLUMN id SET DEFAULT nextval('public.serving_time_id_seq'::regclass);


--
-- Name: servings id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.servings ALTER COLUMN id SET DEFAULT nextval('public.servings_id_seq'::regclass);


--
-- Name: time_slot id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.time_slot ALTER COLUMN id SET DEFAULT nextval('public.time_slot_id_seq'::regclass);


--
-- Name: user_image_uploads id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.user_image_uploads ALTER COLUMN id SET DEFAULT nextval('public.user_image_uploads_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: additives_locales additives_locales_locale_parent_id_unique; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.additives_locales
    ADD CONSTRAINT additives_locales_locale_parent_id_unique UNIQUE (_locale, _parent_id);


--
-- Name: additives_locales additives_locales_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.additives_locales
    ADD CONSTRAINT additives_locales_pkey PRIMARY KEY (id);


--
-- Name: additives additives_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.additives
    ADD CONSTRAINT additives_pkey PRIMARY KEY (id);


--
-- Name: allergens_locales allergens_locales_locale_parent_id_unique; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.allergens_locales
    ADD CONSTRAINT allergens_locales_locale_parent_id_unique UNIQUE (_locale, _parent_id);


--
-- Name: allergens_locales allergens_locales_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.allergens_locales
    ADD CONSTRAINT allergens_locales_pkey PRIMARY KEY (id);


--
-- Name: allergens allergens_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.allergens
    ADD CONSTRAINT allergens_pkey PRIMARY KEY (id);


--
-- Name: features features_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.features
    ADD CONSTRAINT features_pkey PRIMARY KEY (id);


--
-- Name: info info_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.info
    ADD CONSTRAINT info_pkey PRIMARY KEY (id);


--
-- Name: locale locale_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale
    ADD CONSTRAINT locale_pkey PRIMARY KEY (id);


--
-- Name: locale_rels locale_rels_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale_rels
    ADD CONSTRAINT locale_rels_pkey PRIMARY KEY (id);


--
-- Name: media media_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.media
    ADD CONSTRAINT media_pkey PRIMARY KEY (id);


--
-- Name: mensa mensa_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.mensa
    ADD CONSTRAINT mensa_pkey PRIMARY KEY (id);


--
-- Name: mensa_provider mensa_provider_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.mensa_provider
    ADD CONSTRAINT mensa_provider_pkey PRIMARY KEY (id);


--
-- Name: nutrient_labels_locales nutrient_labels_locales_locale_parent_id_unique; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_labels_locales
    ADD CONSTRAINT nutrient_labels_locales_locale_parent_id_unique UNIQUE (_locale, _parent_id);


--
-- Name: nutrient_labels_locales nutrient_labels_locales_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_labels_locales
    ADD CONSTRAINT nutrient_labels_locales_pkey PRIMARY KEY (id);


--
-- Name: nutrient_labels nutrient_labels_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_labels
    ADD CONSTRAINT nutrient_labels_pkey PRIMARY KEY (id);


--
-- Name: nutrient_units nutrient_units_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_units
    ADD CONSTRAINT nutrient_units_pkey PRIMARY KEY (id);


--
-- Name: nutrient_values nutrient_values_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_values
    ADD CONSTRAINT nutrient_values_pkey PRIMARY KEY (id);


--
-- Name: nutrients nutrients_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrients
    ADD CONSTRAINT nutrients_pkey PRIMARY KEY (id);


--
-- Name: payload_locked_documents payload_locked_documents_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents
    ADD CONSTRAINT payload_locked_documents_pkey PRIMARY KEY (id);


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_pkey PRIMARY KEY (id);


--
-- Name: payload_migrations payload_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_migrations
    ADD CONSTRAINT payload_migrations_pkey PRIMARY KEY (id);


--
-- Name: payload_preferences payload_preferences_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_preferences
    ADD CONSTRAINT payload_preferences_pkey PRIMARY KEY (id);


--
-- Name: payload_preferences_rels payload_preferences_rels_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_preferences_rels
    ADD CONSTRAINT payload_preferences_rels_pkey PRIMARY KEY (id);


--
-- Name: recipes recipes_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes
    ADD CONSTRAINT recipes_pkey PRIMARY KEY (id);


--
-- Name: recipes_rels recipes_rels_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes_rels
    ADD CONSTRAINT recipes_rels_pkey PRIMARY KEY (id);


--
-- Name: serving_time serving_time_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.serving_time
    ADD CONSTRAINT serving_time_pkey PRIMARY KEY (id);


--
-- Name: servings servings_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.servings
    ADD CONSTRAINT servings_pkey PRIMARY KEY (id);


--
-- Name: time_slot time_slot_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.time_slot
    ADD CONSTRAINT time_slot_pkey PRIMARY KEY (id);


--
-- Name: user_image_uploads user_image_uploads_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.user_image_uploads
    ADD CONSTRAINT user_image_uploads_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: additives_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX additives_created_at_1_idx ON public.additives USING btree (created_at);


--
-- Name: additives_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX additives_updated_at_1_idx ON public.additives USING btree (updated_at);


--
-- Name: allergens_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX allergens_created_at_1_idx ON public.allergens USING btree (created_at);


--
-- Name: allergens_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX allergens_updated_at_1_idx ON public.allergens USING btree (updated_at);


--
-- Name: features_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX features_created_at_1_idx ON public.features USING btree (created_at);


--
-- Name: features_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX features_updated_at_1_idx ON public.features USING btree (updated_at);


--
-- Name: info_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX info_created_at_1_idx ON public.info USING btree (created_at);


--
-- Name: info_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX info_updated_at_1_idx ON public.info USING btree (updated_at);


--
-- Name: locale_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_created_at_1_idx ON public.locale USING btree (created_at);


--
-- Name: locale_rels_additives_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_rels_additives_id_1_idx ON public.locale_rels USING btree (additives_id);


--
-- Name: locale_rels_allergens_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_rels_allergens_id_1_idx ON public.locale_rels USING btree (allergens_id);


--
-- Name: locale_rels_features_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_rels_features_id_1_idx ON public.locale_rels USING btree (features_id);


--
-- Name: locale_rels_order_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_rels_order_idx ON public.locale_rels USING btree ("order");


--
-- Name: locale_rels_parent_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_rels_parent_idx ON public.locale_rels USING btree (parent_id);


--
-- Name: locale_rels_path_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_rels_path_idx ON public.locale_rels USING btree (path);


--
-- Name: locale_rels_recipes_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_rels_recipes_id_1_idx ON public.locale_rels USING btree (recipes_id);


--
-- Name: locale_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX locale_updated_at_1_idx ON public.locale USING btree (updated_at);


--
-- Name: media_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX media_created_at_1_idx ON public.media USING btree (created_at);


--
-- Name: media_filename_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE UNIQUE INDEX media_filename_1_idx ON public.media USING btree (filename);


--
-- Name: media_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX media_updated_at_1_idx ON public.media USING btree (updated_at);


--
-- Name: mensa_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX mensa_created_at_1_idx ON public.mensa USING btree (created_at);


--
-- Name: mensa_provider_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX mensa_provider_1_idx ON public.mensa USING btree (provider_id);


--
-- Name: mensa_provider_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX mensa_provider_created_at_1_idx ON public.mensa_provider USING btree (created_at);


--
-- Name: mensa_provider_slug_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX mensa_provider_slug_1_idx ON public.mensa_provider USING btree (slug);


--
-- Name: mensa_provider_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX mensa_provider_updated_at_1_idx ON public.mensa_provider USING btree (updated_at);


--
-- Name: mensa_slug_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX mensa_slug_1_idx ON public.mensa USING btree (slug);


--
-- Name: mensa_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX mensa_updated_at_1_idx ON public.mensa USING btree (updated_at);


--
-- Name: nutrient_labels_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrient_labels_created_at_1_idx ON public.nutrient_labels USING btree (created_at);


--
-- Name: nutrient_labels_name_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE UNIQUE INDEX nutrient_labels_name_1_idx ON public.nutrient_labels_locales USING btree (name, _locale);


--
-- Name: nutrient_labels_unit_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrient_labels_unit_1_idx ON public.nutrient_labels USING btree (unit_id);


--
-- Name: nutrient_labels_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrient_labels_updated_at_1_idx ON public.nutrient_labels USING btree (updated_at);


--
-- Name: nutrient_units_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrient_units_created_at_1_idx ON public.nutrient_units USING btree (created_at);


--
-- Name: nutrient_units_name_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE UNIQUE INDEX nutrient_units_name_1_idx ON public.nutrient_units USING btree (name);


--
-- Name: nutrient_units_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrient_units_updated_at_1_idx ON public.nutrient_units USING btree (updated_at);


--
-- Name: nutrient_values_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrient_values_created_at_1_idx ON public.nutrient_values USING btree (created_at);


--
-- Name: nutrient_values_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrient_values_updated_at_1_idx ON public.nutrient_values USING btree (updated_at);


--
-- Name: nutrient_values_value_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE UNIQUE INDEX nutrient_values_value_1_idx ON public.nutrient_values USING btree (value);


--
-- Name: nutrients_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrients_created_at_1_idx ON public.nutrients USING btree (created_at);


--
-- Name: nutrients_nutrient_label_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrients_nutrient_label_1_idx ON public.nutrients USING btree (nutrient_label_id);


--
-- Name: nutrients_nutrient_value_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrients_nutrient_value_1_idx ON public.nutrients USING btree (nutrient_value_id);


--
-- Name: nutrients_recipe_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrients_recipe_1_idx ON public.nutrients USING btree (recipe_id);


--
-- Name: nutrients_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX nutrients_updated_at_1_idx ON public.nutrients USING btree (updated_at);


--
-- Name: payload_locked_documents_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_created_at_1_idx ON public.payload_locked_documents USING btree (created_at);


--
-- Name: payload_locked_documents_global_slug_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_global_slug_1_idx ON public.payload_locked_documents USING btree (global_slug);


--
-- Name: payload_locked_documents_rels_additives_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_additives_id_1_idx ON public.payload_locked_documents_rels USING btree (additives_id);


--
-- Name: payload_locked_documents_rels_allergens_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_allergens_id_1_idx ON public.payload_locked_documents_rels USING btree (allergens_id);


--
-- Name: payload_locked_documents_rels_features_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_features_id_1_idx ON public.payload_locked_documents_rels USING btree (features_id);


--
-- Name: payload_locked_documents_rels_info_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_info_id_1_idx ON public.payload_locked_documents_rels USING btree (info_id);


--
-- Name: payload_locked_documents_rels_locale_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_locale_id_1_idx ON public.payload_locked_documents_rels USING btree (locale_id);


--
-- Name: payload_locked_documents_rels_media_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_media_id_1_idx ON public.payload_locked_documents_rels USING btree (media_id);


--
-- Name: payload_locked_documents_rels_mensa_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_mensa_id_1_idx ON public.payload_locked_documents_rels USING btree (mensa_id);


--
-- Name: payload_locked_documents_rels_mensa_provider_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_mensa_provider_id_1_idx ON public.payload_locked_documents_rels USING btree (mensa_provider_id);


--
-- Name: payload_locked_documents_rels_nutrient_labels_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_nutrient_labels_id_1_idx ON public.payload_locked_documents_rels USING btree (nutrient_labels_id);


--
-- Name: payload_locked_documents_rels_nutrient_units_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_nutrient_units_id_1_idx ON public.payload_locked_documents_rels USING btree (nutrient_units_id);


--
-- Name: payload_locked_documents_rels_nutrient_values_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_nutrient_values_id_1_idx ON public.payload_locked_documents_rels USING btree (nutrient_values_id);


--
-- Name: payload_locked_documents_rels_nutrients_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_nutrients_id_1_idx ON public.payload_locked_documents_rels USING btree (nutrients_id);


--
-- Name: payload_locked_documents_rels_order_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_order_idx ON public.payload_locked_documents_rels USING btree ("order");


--
-- Name: payload_locked_documents_rels_parent_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_parent_idx ON public.payload_locked_documents_rels USING btree (parent_id);


--
-- Name: payload_locked_documents_rels_path_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_path_idx ON public.payload_locked_documents_rels USING btree (path);


--
-- Name: payload_locked_documents_rels_recipes_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_recipes_id_1_idx ON public.payload_locked_documents_rels USING btree (recipes_id);


--
-- Name: payload_locked_documents_rels_serving_time_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_serving_time_id_1_idx ON public.payload_locked_documents_rels USING btree (serving_time_id);


--
-- Name: payload_locked_documents_rels_servings_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_servings_id_1_idx ON public.payload_locked_documents_rels USING btree (servings_id);


--
-- Name: payload_locked_documents_rels_time_slot_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_time_slot_id_1_idx ON public.payload_locked_documents_rels USING btree (time_slot_id);


--
-- Name: payload_locked_documents_rels_user_image_uploads_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_user_image_uploads_id_1_idx ON public.payload_locked_documents_rels USING btree (user_image_uploads_id);


--
-- Name: payload_locked_documents_rels_users_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_rels_users_id_1_idx ON public.payload_locked_documents_rels USING btree (users_id);


--
-- Name: payload_locked_documents_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_locked_documents_updated_at_1_idx ON public.payload_locked_documents USING btree (updated_at);


--
-- Name: payload_migrations_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_migrations_created_at_1_idx ON public.payload_migrations USING btree (created_at);


--
-- Name: payload_migrations_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_migrations_updated_at_1_idx ON public.payload_migrations USING btree (updated_at);


--
-- Name: payload_preferences_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_preferences_created_at_1_idx ON public.payload_preferences USING btree (created_at);


--
-- Name: payload_preferences_key_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_preferences_key_1_idx ON public.payload_preferences USING btree (key);


--
-- Name: payload_preferences_rels_order_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_preferences_rels_order_idx ON public.payload_preferences_rels USING btree ("order");


--
-- Name: payload_preferences_rels_parent_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_preferences_rels_parent_idx ON public.payload_preferences_rels USING btree (parent_id);


--
-- Name: payload_preferences_rels_path_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_preferences_rels_path_idx ON public.payload_preferences_rels USING btree (path);


--
-- Name: payload_preferences_rels_users_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_preferences_rels_users_id_1_idx ON public.payload_preferences_rels USING btree (users_id);


--
-- Name: payload_preferences_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX payload_preferences_updated_at_1_idx ON public.payload_preferences USING btree (updated_at);


--
-- Name: recipes_ai_thumbnail_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_ai_thumbnail_1_idx ON public.recipes USING btree (ai_thumbnail_id);


--
-- Name: recipes_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_created_at_1_idx ON public.recipes USING btree (created_at);


--
-- Name: recipes_mensa_provider_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_mensa_provider_1_idx ON public.recipes USING btree (mensa_provider_id);


--
-- Name: recipes_rels_additives_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_rels_additives_id_1_idx ON public.recipes_rels USING btree (additives_id);


--
-- Name: recipes_rels_allergens_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_rels_allergens_id_1_idx ON public.recipes_rels USING btree (allergens_id);


--
-- Name: recipes_rels_features_id_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_rels_features_id_1_idx ON public.recipes_rels USING btree (features_id);


--
-- Name: recipes_rels_order_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_rels_order_idx ON public.recipes_rels USING btree ("order");


--
-- Name: recipes_rels_parent_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_rels_parent_idx ON public.recipes_rels USING btree (parent_id);


--
-- Name: recipes_rels_path_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_rels_path_idx ON public.recipes_rels USING btree (path);


--
-- Name: recipes_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX recipes_updated_at_1_idx ON public.recipes USING btree (updated_at);


--
-- Name: serving_time_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX serving_time_created_at_1_idx ON public.serving_time USING btree (created_at);


--
-- Name: serving_time_mensa_info_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX serving_time_mensa_info_1_idx ON public.serving_time USING btree (mensa_info_id);


--
-- Name: serving_time_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX serving_time_updated_at_1_idx ON public.serving_time USING btree (updated_at);


--
-- Name: servings_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX servings_created_at_1_idx ON public.servings USING btree (created_at);


--
-- Name: servings_mensa_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX servings_mensa_1_idx ON public.servings USING btree (mensa_id);


--
-- Name: servings_recipe_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX servings_recipe_1_idx ON public.servings USING btree (recipe_id);


--
-- Name: servings_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX servings_updated_at_1_idx ON public.servings USING btree (updated_at);


--
-- Name: time_slot_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX time_slot_created_at_1_idx ON public.time_slot USING btree (created_at);


--
-- Name: time_slot_serving_time_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX time_slot_serving_time_1_idx ON public.time_slot USING btree (serving_time_id);


--
-- Name: time_slot_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX time_slot_updated_at_1_idx ON public.time_slot USING btree (updated_at);


--
-- Name: user_image_uploads_image_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX user_image_uploads_image_1_idx ON public.user_image_uploads USING btree (image_id);


--
-- Name: user_image_uploads_recipe_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX user_image_uploads_recipe_1_idx ON public.user_image_uploads USING btree (recipe_id);


--
-- Name: user_image_uploads_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX user_image_uploads_updated_at_1_idx ON public.user_image_uploads USING btree (updated_at);


--
-- Name: users_created_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX users_created_at_1_idx ON public.users USING btree (created_at);


--
-- Name: users_email_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE UNIQUE INDEX users_email_1_idx ON public.users USING btree (email);


--
-- Name: users_updated_at_1_idx; Type: INDEX; Schema: public; Owner: mensauser
--

CREATE INDEX users_updated_at_1_idx ON public.users USING btree (updated_at);


--
-- Name: additives_locales additives_locales_parent_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.additives_locales
    ADD CONSTRAINT additives_locales_parent_id_fk FOREIGN KEY (_parent_id) REFERENCES public.additives(id) ON DELETE CASCADE;


--
-- Name: allergens_locales allergens_locales_parent_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.allergens_locales
    ADD CONSTRAINT allergens_locales_parent_id_fk FOREIGN KEY (_parent_id) REFERENCES public.allergens(id) ON DELETE CASCADE;


--
-- Name: locale_rels locale_rels_additives_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale_rels
    ADD CONSTRAINT locale_rels_additives_fk FOREIGN KEY (additives_id) REFERENCES public.additives(id) ON DELETE CASCADE;


--
-- Name: locale_rels locale_rels_allergens_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale_rels
    ADD CONSTRAINT locale_rels_allergens_fk FOREIGN KEY (allergens_id) REFERENCES public.allergens(id) ON DELETE CASCADE;


--
-- Name: locale_rels locale_rels_features_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale_rels
    ADD CONSTRAINT locale_rels_features_fk FOREIGN KEY (features_id) REFERENCES public.features(id) ON DELETE CASCADE;


--
-- Name: locale_rels locale_rels_parent_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale_rels
    ADD CONSTRAINT locale_rels_parent_fk FOREIGN KEY (parent_id) REFERENCES public.locale(id) ON DELETE CASCADE;


--
-- Name: locale_rels locale_rels_recipes_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.locale_rels
    ADD CONSTRAINT locale_rels_recipes_fk FOREIGN KEY (recipes_id) REFERENCES public.recipes(id) ON DELETE CASCADE;


--
-- Name: mensa mensa_provider_id_mensa_provider_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.mensa
    ADD CONSTRAINT mensa_provider_id_mensa_provider_id_fk FOREIGN KEY (provider_id) REFERENCES public.mensa_provider(id) ON DELETE SET NULL;


--
-- Name: nutrient_labels_locales nutrient_labels_locales_parent_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_labels_locales
    ADD CONSTRAINT nutrient_labels_locales_parent_id_fk FOREIGN KEY (_parent_id) REFERENCES public.nutrient_labels(id) ON DELETE CASCADE;


--
-- Name: nutrient_labels nutrient_labels_unit_id_nutrient_units_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrient_labels
    ADD CONSTRAINT nutrient_labels_unit_id_nutrient_units_id_fk FOREIGN KEY (unit_id) REFERENCES public.nutrient_units(id) ON DELETE SET NULL;


--
-- Name: nutrients nutrients_nutrient_label_id_nutrient_labels_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrients
    ADD CONSTRAINT nutrients_nutrient_label_id_nutrient_labels_id_fk FOREIGN KEY (nutrient_label_id) REFERENCES public.nutrient_labels(id) ON DELETE SET NULL;


--
-- Name: nutrients nutrients_nutrient_value_id_nutrient_values_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrients
    ADD CONSTRAINT nutrients_nutrient_value_id_nutrient_values_id_fk FOREIGN KEY (nutrient_value_id) REFERENCES public.nutrient_values(id) ON DELETE SET NULL;


--
-- Name: nutrients nutrients_recipe_id_recipes_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.nutrients
    ADD CONSTRAINT nutrients_recipe_id_recipes_id_fk FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON DELETE SET NULL;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_additives_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_additives_fk FOREIGN KEY (additives_id) REFERENCES public.additives(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_allergens_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_allergens_fk FOREIGN KEY (allergens_id) REFERENCES public.allergens(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_features_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_features_fk FOREIGN KEY (features_id) REFERENCES public.features(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_info_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_info_fk FOREIGN KEY (info_id) REFERENCES public.info(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_locale_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_locale_fk FOREIGN KEY (locale_id) REFERENCES public.locale(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_media_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_media_fk FOREIGN KEY (media_id) REFERENCES public.media(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_mensa_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_mensa_fk FOREIGN KEY (mensa_id) REFERENCES public.mensa(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_mensa_provider_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_mensa_provider_fk FOREIGN KEY (mensa_provider_id) REFERENCES public.mensa_provider(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_nutrient_labels_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_nutrient_labels_fk FOREIGN KEY (nutrient_labels_id) REFERENCES public.nutrient_labels(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_nutrient_units_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_nutrient_units_fk FOREIGN KEY (nutrient_units_id) REFERENCES public.nutrient_units(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_nutrient_values_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_nutrient_values_fk FOREIGN KEY (nutrient_values_id) REFERENCES public.nutrient_values(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_nutrients_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_nutrients_fk FOREIGN KEY (nutrients_id) REFERENCES public.nutrients(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_parent_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_parent_fk FOREIGN KEY (parent_id) REFERENCES public.payload_locked_documents(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_recipes_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_recipes_fk FOREIGN KEY (recipes_id) REFERENCES public.recipes(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_serving_time_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_serving_time_fk FOREIGN KEY (serving_time_id) REFERENCES public.serving_time(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_servings_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_servings_fk FOREIGN KEY (servings_id) REFERENCES public.servings(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_time_slot_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_time_slot_fk FOREIGN KEY (time_slot_id) REFERENCES public.time_slot(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_user_image_uploads_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_user_image_uploads_fk FOREIGN KEY (user_image_uploads_id) REFERENCES public.user_image_uploads(id) ON DELETE CASCADE;


--
-- Name: payload_locked_documents_rels payload_locked_documents_rels_users_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_locked_documents_rels
    ADD CONSTRAINT payload_locked_documents_rels_users_fk FOREIGN KEY (users_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: payload_preferences_rels payload_preferences_rels_parent_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_preferences_rels
    ADD CONSTRAINT payload_preferences_rels_parent_fk FOREIGN KEY (parent_id) REFERENCES public.payload_preferences(id) ON DELETE CASCADE;


--
-- Name: payload_preferences_rels payload_preferences_rels_users_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.payload_preferences_rels
    ADD CONSTRAINT payload_preferences_rels_users_fk FOREIGN KEY (users_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: recipes recipes_ai_thumbnail_id_media_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes
    ADD CONSTRAINT recipes_ai_thumbnail_id_media_id_fk FOREIGN KEY (ai_thumbnail_id) REFERENCES public.media(id) ON DELETE SET NULL;


--
-- Name: recipes recipes_mensa_provider_id_mensa_provider_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes
    ADD CONSTRAINT recipes_mensa_provider_id_mensa_provider_id_fk FOREIGN KEY (mensa_provider_id) REFERENCES public.mensa_provider(id) ON DELETE SET NULL;


--
-- Name: recipes_rels recipes_rels_additives_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes_rels
    ADD CONSTRAINT recipes_rels_additives_fk FOREIGN KEY (additives_id) REFERENCES public.additives(id) ON DELETE CASCADE;


--
-- Name: recipes_rels recipes_rels_allergens_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes_rels
    ADD CONSTRAINT recipes_rels_allergens_fk FOREIGN KEY (allergens_id) REFERENCES public.allergens(id) ON DELETE CASCADE;


--
-- Name: recipes_rels recipes_rels_features_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes_rels
    ADD CONSTRAINT recipes_rels_features_fk FOREIGN KEY (features_id) REFERENCES public.features(id) ON DELETE CASCADE;


--
-- Name: recipes_rels recipes_rels_parent_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.recipes_rels
    ADD CONSTRAINT recipes_rels_parent_fk FOREIGN KEY (parent_id) REFERENCES public.recipes(id) ON DELETE CASCADE;


--
-- Name: serving_time serving_time_mensa_info_id_mensa_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.serving_time
    ADD CONSTRAINT serving_time_mensa_info_id_mensa_id_fk FOREIGN KEY (mensa_info_id) REFERENCES public.mensa(id) ON DELETE SET NULL;


--
-- Name: servings servings_mensa_id_mensa_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.servings
    ADD CONSTRAINT servings_mensa_id_mensa_id_fk FOREIGN KEY (mensa_id) REFERENCES public.mensa(id) ON DELETE SET NULL;


--
-- Name: servings servings_recipe_id_recipes_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.servings
    ADD CONSTRAINT servings_recipe_id_recipes_id_fk FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON DELETE SET NULL;


--
-- Name: time_slot time_slot_serving_time_id_serving_time_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.time_slot
    ADD CONSTRAINT time_slot_serving_time_id_serving_time_id_fk FOREIGN KEY (serving_time_id) REFERENCES public.serving_time(id) ON DELETE SET NULL;


--
-- Name: user_image_uploads user_image_uploads_image_id_media_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.user_image_uploads
    ADD CONSTRAINT user_image_uploads_image_id_media_id_fk FOREIGN KEY (image_id) REFERENCES public.media(id) ON DELETE SET NULL;


--
-- Name: user_image_uploads user_image_uploads_recipe_id_recipes_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: mensauser
--

ALTER TABLE ONLY public.user_image_uploads
    ADD CONSTRAINT user_image_uploads_recipe_id_recipes_id_fk FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

