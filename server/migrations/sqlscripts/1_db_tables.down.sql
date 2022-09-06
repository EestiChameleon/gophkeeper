BEGIN;
------------
-- TABLES --
------------
DROP TABLE IF EXISTS gophkeeper_users;
DROP TABLE IF EXISTS gk_pair;
DROP TABLE IF EXISTS gk_text;
DROP TABLE IF EXISTS gk_bin;
DROP TABLE IF EXISTS gk_card;

DROP INDEX IF EXISTS gk_pair_user_id_title_version_uindex;
DROP INDEX IF EXISTS gk_text_user_id_title_version_uindex;
DROP INDEX IF EXISTS gk_bin_user_id_title_version_uindex;
DROP INDEX IF EXISTS gk_card_user_id_title_version_uindex;

----------
-- DATA --
----------

COMMIT;