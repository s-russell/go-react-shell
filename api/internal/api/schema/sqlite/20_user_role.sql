create table
    if not exists user_role (
        id integer primary KEY,
        user_id integer,
        role_id integer,
        FOREIGN key (user_id) REFERENCES user (id),
        FOREIGN key (role_id) REFERENCES roles (id)
    );