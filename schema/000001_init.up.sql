--migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
create table expressions
(
    id               serial                  not null unique,
    value            varchar(255)            not null,
    result           integer,
    state            varchar(255)            not null,
    created_at       timestamp default now() not null,
    eval_started_at  timestamp,
    eval_finished_at timestamp
);