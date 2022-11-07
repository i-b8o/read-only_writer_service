BEGIN;

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;

SET SCHEMA 'public';

-- CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- TABLES --
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS speech;
DROP MATERIALIZED VIEW IF EXISTS reg_search;
DROP INDEX IF EXISTS idx_search;
DROP TABLE IF EXISTS paragraphs;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS regulations;


CREATE TABLE regulations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (NAME != '') UNIQUE,
    abbreviation TEXT,
    title TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE regulations ADD COLUMN ts tsvector GENERATED ALWAYS AS (setweight(to_tsvector('russian', coalesce(name, '')), 'A') || setweight(to_tsvector('russian', coalesce(title, '')), 'B')) STORED;
CREATE INDEX reg_ts_idx ON regulations USING GIN (ts);



CREATE TABLE chapters (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (name != ''),
    order_num SMALLINT NOT NULL CHECK (order_num >= 0),
    num TEXT,
    r_id integer REFERENCES regulations,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE chapters ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('russian', name)) STORED;
CREATE INDEX ch_ts_idx ON chapters USING GIN (ts);

CREATE TABLE paragraphs (
    id SERIAL PRIMARY KEY,
    paragraph_id INT NOT NULL CHECK (paragraph_id >= 0),
    order_num INT NOT NULL CHECK (order_num >= 0),
    is_table BOOLEAN NOT NULL,
    is_nft BOOLEAN NOT NULL,
    has_links BOOLEAN NOT NULL,
    class TEXT,
    content TEXT NOT NULL,
    c_id integer REFERENCES chapters
);

ALTER TABLE paragraphs ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('russian', content)) STORED;
CREATE INDEX p_ts_idx ON paragraphs USING GIN (ts);


CREATE MATERIALIZED VIEW reg_search 
AS SELECT 
r.id AS "r_id", r.name AS "r_name", NULL AS "c_id", NULL AS "c_name", CAST(NULL AS integer) AS "p_id", NULL AS "p_text", r.name AS "text",
to_tsvector('russian', r.name) AS ts FROM regulations AS r UNION SELECT 
NULL AS "r_id", r.name AS "r_name", c.id AS "c_id", c.name AS "c_name", NULL AS "p_id", NULL AS "p_text", c.name AS "text",
to_tsvector('russian', c.name) AS ts FROM chapters AS c INNER JOIN regulations AS r ON r.id= c.r_id
UNION SELECT 
NULL AS "r_id", r.name AS "r_name", c.id AS "c_id", c.name AS "c_name", p.paragraph_id AS "p_id", p.content AS "p_text", p.content AS "text",
to_tsvector('russian', content) AS ts 
FROM paragraphs AS p INNER JOIN chapters AS c ON p.c_id= c.id INNER JOIN regulations AS r ON c.r_id = r.id;

create index idx_search on reg_search using GIN(ts);