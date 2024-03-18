create procedure insert_all_actors(f_id int, actors_id int[]) language plpgsql
as $$
declare
a_id int;
begin
        foreach a_id in array actors_id
            loop
                insert into actors_has_films (film_id, actor_id) values (f_id, a_id);
end loop;
end;
$$