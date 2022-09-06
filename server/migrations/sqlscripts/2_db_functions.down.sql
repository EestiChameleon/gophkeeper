BEGIN;

---------------
-- FUNCTIONS --
---------------
DROP FUNCTION IF EXISTS users_list();
DROP FUNCTION IF EXISTS user_by_login(f_login varchar);
DROP FUNCTION IF EXISTS users_add(f_login varchar, f_password varchar);

DROP FUNCTION IF EXISTS pair_by_title(f_title varchar, f_user_id int);
DROP FUNCTION IF EXISTS pair_add(f_user_id int, f_title varchar, f_login varchar, f_pass varchar, f_comment varchar, f_version smallint);
DROP FUNCTION IF EXISTS pair_del_by_title(f_title varchar, f_user_id int);

DROP FUNCTION IF EXISTS text_by_title(f_title varchar, f_user_id int);
DROP FUNCTION IF EXISTS text_add(f_user_id int, f_title varchar, f_body varchar, f_comment varchar, f_version smallint);
DROP FUNCTION IF EXISTS text_del_by_title(f_title varchar, f_user_id int);

DROP FUNCTION IF EXISTS bin_by_title(f_title varchar, f_user_id int);
DROP FUNCTION IF EXISTS bin_add(f_user_id int, f_title varchar, f_body bytea, f_comment varchar, f_version smallint);
DROP FUNCTION IF EXISTS bin_del_by_title(f_title varchar, f_user_id int);

DROP FUNCTION IF EXISTS card_by_title(f_title varchar, f_user_id int);
DROP FUNCTION IF EXISTS card_add(f_user_id int, f_title varchar, f_number varchar, f_expdate varchar, f_comment varchar, f_version smallint);
DROP FUNCTION IF EXISTS card_del_by_title(f_title varchar, f_user_id int);

DROP FUNCTION IF EXISTS pairs_all_last_version_by_user_id(f_user_id int);
DROP FUNCTION IF EXISTS texts_all_last_version_by_user_id(f_user_id int);
DROP FUNCTION IF EXISTS bins_all_last_version_by_user_id(f_user_id int);
DROP FUNCTION IF EXISTS cards_all_last_version_by_user_id(f_user_id int);


COMMIT;