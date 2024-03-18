create function is_exist(l text, pass text) returns bool
    language plpgsql
as $$
begin
return exists(select * from login_forms where login = 'Zuzinho' and password = '12345');
end;
$$