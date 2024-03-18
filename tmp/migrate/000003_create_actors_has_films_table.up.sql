create table actors_has_films (
                                  film_id int not null,
                                  actor_id int not null,
                                  foreign key(film_id) references films(film_id)
                                      on delete cascade
                                      on update cascade,
                                  foreign key(actor_id) references actors(actor_id)
                                      on delete cascade
                                      on update cascade
);