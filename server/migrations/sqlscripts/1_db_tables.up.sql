BEGIN;
------------
-- TABLES --
------------
-- // using create sequence instead of serial type

CREATE TABLE IF NOT EXISTS gophkeeper_users
(
    id       serial primary key,
    login    varchar not null unique,
    password varchar not null
);

CREATE TABLE IF NOT EXISTS gk_pair
(
    id         serial primary key,
    user_id    int                not null,
    title      varchar            not null,
    login      varchar            not null,
    pass       varchar            not null,
    comment    varchar,
    version    smallint default 1 not null,
    deleted_at timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS gk_pair_user_id_title_version_uindex
    on gk_pair (user_id, title, version);

CREATE TABLE IF NOT EXISTS gk_text
(
    id         serial primary key,
    user_id    int                not null,
    title      varchar            not null,
    body       varchar            not null,
    comment    varchar,
    version    smallint default 1 not null,
    deleted_at timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS gk_text_user_id_title_version_uindex
    on gk_text (user_id, title, version);

CREATE TABLE IF NOT EXISTS gk_bin
(
    id         serial primary key,
    user_id    int                not null,
    title      varchar            not null,
    body       bytea              not null,
    comment    varchar,
    version    smallint default 1 not null,
    deleted_at timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS gk_bin_user_id_title_version_uindex
    on gk_bin (user_id, title, version);

CREATE TABLE IF NOT EXISTS gk_card
(
    id         serial primary key,
    user_id    int                not null,
    title      varchar            not null,
    number     varchar(16)        not null,
    expdate    varchar(12)        not null,
    comment    varchar,
    version    smallint default 1 not null ,
    deleted_at timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS gk_card_user_id_title_version_uindex
    on gk_card (user_id, title, version);


----------
-- DATA --
----------


COMMIT;