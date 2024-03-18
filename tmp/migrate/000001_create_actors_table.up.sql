create type gender_enum as enum ('Male', 'Female');

create table actors (
                        actor_id serial not null primary key,
                        first_name text not null,
                        second_name text not null,
                        gender gender_enum not null,
                        birthday date not null
);