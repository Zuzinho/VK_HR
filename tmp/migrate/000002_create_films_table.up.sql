create table films (
                       film_id serial not null primary key,
                       name text not null,
                       description text not null,
                       premier_date date not null,
                       rating float4 not null
);