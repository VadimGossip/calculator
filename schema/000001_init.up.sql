--migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/calculator_db?sslmode=disable" up
create table expressions
(
    id               serial                  not null unique,
    value            varchar(255)            not null,
    result           numeric(10, 5),
    state            varchar(255)            not null,
    error_msg        varchar(2000),
    created_at       timestamp default now() not null,
    eval_started_at  timestamp,
    eval_finished_at timestamp
);

create table agents
(
    name              varchar(255)  not null unique,
    created_at        timestamp default now() not null,
    last_heartbeat_at timestamp default now() not null
);

create table operation_durations
(
    operation_name   varchar(255)  not null unique,
    duration         numeric(7) not null,
    created_at       timestamp default now() not null,
    updated_at       timestamp default now() not null
);

create table sub_expressions
(
    id                 serial                  not null unique,
    expression_id      integer,
    val1               numeric(10, 5),
    val2               numeric(10, 5),
    sub_expression_id1 integer,
    sub_expression_id2 integer,
    operation_name     varchar(255) not null,
    result             numeric(10, 5),
    agent_name         varchar(255),
    eval_started_at    timestamp,
    eval_finished_at   timestamp,
    is_last            boolean default false
);
