create table
    if not exists user (
        id integer primary key,
        first_name text,
        last_name text,
        username text,
        password text
    );