--migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/calculator_db?sslmode=disable" up
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

create table agents
(
    name   varchar(255)  not null unique,
    created_at       timestamp default now() not null,
    last_heartbeat_at timestamp default now() not null
);

create table operation_durations
(
    operation_name   varchar(255)  not null unique,
    duration         numeric(3) not null,
    created_at       timestamp default now() not null,
    updated_at       timestamp default now() not null
)

