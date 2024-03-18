create type role_enum as enum ('Regular User', 'Admin');

create table users (
                       login text not null primary key,
                       role role_enum not null default 'Regular User'
);