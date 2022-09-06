BEGIN;

---------------
-- FUNCTIONS --
---------------
-- users List
CREATE OR REPLACE FUNCTION users_list()
    RETURNS TABLE (
                      id int,
                      login varchar,
                      password varchar
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT  u.id, u.login, u.password
        FROM gophkeeper_users AS u
        ORDER BY u.id;
END
$$ LANGUAGE plpgsql;

-- SELECT By login
CREATE OR REPLACE FUNCTION user_by_login(f_login varchar)
    RETURNS TABLE (
                      id int,
                      login varchar,
                      password varchar
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT  u.id, u.login, u.password
        FROM gophkeeper_users AS u
        WHERE u.login = f_login;
END
$$ LANGUAGE plpgsql;

-- users ADD insert new row
CREATE OR REPLACE FUNCTION users_add(f_login varchar, f_password varchar)
    RETURNS integer AS
$result$
DECLARE
    result integer;
BEGIN
    INSERT INTO gophkeeper_users (login, password)
    VALUES (f_login, f_password)
    RETURNING id INTO result;
    RETURN result;
END;
$result$
    LANGUAGE plpgsql;


-- gk_pair table
-- SELECT by title last version gk_pair row for user id
CREATE OR REPLACE FUNCTION pair_by_title(f_title varchar, f_user_id int)
    RETURNS TABLE (
                      id            int,
                      user_id int,
                      title     varchar,
                      login     varchar,
                      pass    varchar,
                      comment varchar,
                      version smallint,
                      deleted_at timestamp

                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT  pd.id, pd.user_id, pd.title, pd.login, pd.pass, pd.comment, pd.version, pd.deleted_at
        FROM gk_pair AS pd
        WHERE pd.title = f_title and pd.user_id = f_user_id
        ORDER BY pd.version DESC LIMIT 1;
END
$$ LANGUAGE plpgsql;

-- gk_pair ADD insert new row
CREATE OR REPLACE FUNCTION pair_add(
    f_user_id       int,
    f_title          varchar,
    f_login varchar,
    f_pass varchar,
    f_comment varchar,
    f_version smallint
)
    RETURNS integer AS
$result$
DECLARE
    result integer;
BEGIN
    INSERT INTO gk_pair (user_id, title, login, pass, comment, version)
    VALUES (f_user_id, f_title, f_login, f_pass, f_comment, f_version)
    RETURNING id INTO result;
    RETURN result;
END;
$result$
    LANGUAGE plpgsql;

-- gk_pair DELETE by title
CREATE OR REPLACE FUNCTION pair_del_by_title(f_title varchar, f_user_id int)
    RETURNS int AS $aff_rows$
DECLARE aff_rows int;
BEGIN
    UPDATE gk_pair AS pd
    SET deleted_at = current_timestamp
    WHERE pd.title = f_title AND pd.user_id = f_user_id;
    GET DIAGNOSTICS aff_rows = ROW_COUNT;
    RETURN aff_rows;
END;
$aff_rows$ LANGUAGE plpgsql;


-- gk_text table
-- SELECT by title last version gk_text row for user id
CREATE OR REPLACE FUNCTION text_by_title(f_title varchar, f_user_id int)
    RETURNS TABLE (
                      id            int,
                      user_id int,
                      title     varchar,
                      body   varchar,
                      comment varchar,
                      version smallint,
                      deleted_at timestamp

                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT  td.id, td.user_id, td.title, td.body, td.comment, td.version, td.deleted_at
        FROM gk_text AS td
        WHERE td.title = f_title and td.user_id = f_user_id
        ORDER BY td.version DESC LIMIT 1;
END
$$ LANGUAGE plpgsql;

-- gk_text ADD insert new row
CREATE OR REPLACE FUNCTION text_add(
    f_user_id       int,
    f_title          varchar,
    f_body varchar,
    f_comment varchar,
    f_version smallint
)
    RETURNS integer AS
$result$
DECLARE
    result integer;
BEGIN
    INSERT INTO gk_text (user_id, title, body, comment, version)
    VALUES (f_user_id, f_title, f_body, f_comment, f_version)
    RETURNING id INTO result;
    RETURN result;
END;
$result$
    LANGUAGE plpgsql;

-- gk_text DELETE by title sets the deleted_at value equal to current timestamp
CREATE OR REPLACE FUNCTION text_del_by_title(f_title varchar, f_user_id int)
    RETURNS int AS $aff_rows$
DECLARE aff_rows int;
BEGIN
    UPDATE gk_text AS td
    SET deleted_at = current_timestamp
    WHERE td.title = f_title AND user_id = f_user_id;
    GET DIAGNOSTICS aff_rows = ROW_COUNT;
    RETURN aff_rows;
END;
$aff_rows$ LANGUAGE plpgsql;


-- gk_bin table
-- SELECT by title last version gk_bin row for user id
CREATE OR REPLACE FUNCTION bin_by_title(f_title varchar, f_user_id int)
    RETURNS TABLE (
                      id            int,
                      user_id int ,
                      title     varchar,
                      body     bytea,
                      comment varchar,
                      version smallint,
                      deleted_at timestamp
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT  bd.id, bd.user_id, bd.title, bd.body, bd.comment, bd.version, bd.deleted_at
        FROM gk_bin AS bd
        WHERE bd.title = f_title and bd.user_id = f_user_id
        ORDER BY bd.version DESC LIMIT 1;
END
$$ LANGUAGE plpgsql;

-- gk_bin ADD insert new row
CREATE OR REPLACE FUNCTION bin_add(
    f_user_id       int,
    f_title          varchar,
    f_body bytea,
    f_comment varchar,
    f_version smallint
)
    RETURNS integer AS
$result$
DECLARE
    result integer;
BEGIN
    INSERT INTO gk_bin (user_id, title, body, comment, version)
    VALUES (f_user_id, f_title, f_body, f_comment, f_version)
    RETURNING id INTO result;
    RETURN result;
END;
$result$
    LANGUAGE plpgsql;

-- gk_bin DELETE by title sets the deleted_at value equal to current timestamp
CREATE OR REPLACE FUNCTION bin_del_by_title(f_title varchar, f_user_id int)
    RETURNS int AS $aff_rows$
DECLARE aff_rows int;
BEGIN
    UPDATE gk_bin AS bd
    SET deleted_at = current_timestamp
    WHERE bd.title = f_title AND bd.user_id = f_user_id;
    GET DIAGNOSTICS aff_rows = ROW_COUNT;
    RETURN aff_rows;
END;
$aff_rows$ LANGUAGE plpgsql;


-- card table
-- SELECT by title last version carddate row for user id
CREATE OR REPLACE FUNCTION card_by_title(f_title varchar, f_user_id int)
    RETURNS TABLE (
                      id            int,
                      user_id int ,
                      title     varchar,
                      number     varchar,
                      expdate varchar,
                      comment varchar,
                      version smallint,
                      deleted_at timestamp
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT cd.id, cd.user_id, cd.title, cd.number, cd.expdate, cd.comment, cd.version, cd.deleted_at
        FROM gk_card AS cd
        WHERE cd.title = f_title and cd.user_id = f_user_id
        ORDER BY cd.version DESC LIMIT 1;
END
$$ LANGUAGE plpgsql;

-- gk_card ADD insert new row
CREATE OR REPLACE FUNCTION card_add(
    f_user_id       int,
    f_title          varchar,
    f_number varchar,
    f_expdate varchar,
    f_comment varchar,
    f_version smallint
)
    RETURNS integer AS
$result$
DECLARE
    result integer;
BEGIN
    INSERT INTO gk_card (user_id, title, number, expdate, comment, version)
    VALUES (f_user_id, f_title, f_number, f_expdate, f_comment, f_version)
    RETURNING id INTO result;
    RETURN result;
END;
$result$
    LANGUAGE plpgsql;

-- gk_card DELETE by title sets the deleted_at value equal to current timestamp
CREATE OR REPLACE FUNCTION card_del_by_title(f_title varchar, f_user_id int)
    RETURNS int AS $aff_rows$
DECLARE aff_rows int;
BEGIN
    UPDATE gk_card AS cd
    SET deleted_at = current_timestamp
    WHERE cd.title = f_title AND cd.user_id = f_user_id;
    GET DIAGNOSTICS aff_rows = ROW_COUNT;
    RETURN aff_rows;
END;
$aff_rows$ LANGUAGE plpgsql;


-- find users last version data for sync
-- pair data
CREATE OR REPLACE FUNCTION pairs_all_last_version_by_user_id(f_user_id int)
    RETURNS TABLE (
                      title     varchar,
                      login     varchar,
                      pass    varchar,
                      comment varchar,
                      version smallint
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT DISTINCT ON (pd.title) pd.title, pd.login, pd.pass, pd.comment, pd.version
        FROM gk_pair AS pd
        where pd.user_id = f_user_id AND deleted_at isnull
        order by pd.title, pd.version desc;
END
$$ LANGUAGE plpgsql;

-- text data
CREATE OR REPLACE FUNCTION texts_all_last_version_by_user_id(f_user_id int)
    RETURNS TABLE (
                      title     varchar,
                      body   varchar,
                      comment varchar,
                      version smallint
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT DISTINCT ON (t.title) t.title, t.body, t.comment, t.version
        FROM gk_text AS t
        WHERE t.user_id = f_user_id AND deleted_at isnull
        ORDER BY t.title, t.version DESC;
END
$$ LANGUAGE plpgsql;

-- bin data
CREATE OR REPLACE FUNCTION bins_all_last_version_by_user_id(f_user_id int)
    RETURNS TABLE (
                      title     varchar,
                      body     bytea,
                      comment varchar,
                      version smallint
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT DISTINCT ON (b.title) b.title, b.body, b.comment, b.version
        FROM gk_bin AS b
        WHERE b.user_id = f_user_id AND deleted_at isnull
        ORDER BY b.title, b.version DESC;
END
$$ LANGUAGE plpgsql;

-- pair data
CREATE OR REPLACE FUNCTION cards_all_last_version_by_user_id(f_user_id int)
    RETURNS TABLE (
                      title     varchar,
                      number     varchar,
                      expdate varchar,
                      comment varchar,
                      version smallint
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT DISTINCT ON (cd.title) cd.title, cd.number, cd.expdate, cd.comment, cd.version
        FROM gk_card AS cd
        WHERE cd.user_id = f_user_id AND deleted_at isnull
        ORDER BY cd.title, cd.version DESC;
END
$$ LANGUAGE plpgsql;

COMMIT;