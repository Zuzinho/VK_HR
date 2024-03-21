create table login_forms (
                             login text not null primary key,
                             password text not null
);

insert into login_forms (login, password) values ('Zuzinho', '12345');