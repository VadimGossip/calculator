--migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/calculator_db?sslmode=disable" up

create table users (
   id            serial not null primary key,
   login         varchar(255) not null unique,
   password      varchar(255) not null,
   admin         boolean default false,
   registered_at timestamp default now() not null
);

create table refresh_tokens(
  id            serial not null unique,
  user_id       int  not null,
  token         varchar(255) not null,
  expires_at    timestamp not null
);

alter table refresh_tokens add foreign key (user_id) REFERENCES users;
create index in_refresh_tokens_user_id on refresh_tokens(user_id);

create table expressions
(
    id               serial    primary key,
    user_id          int              not null,
    req_uid          varchar(2000)           not null,
    value            varchar(255)            not null,
    result           numeric(50, 5),
    state            varchar(255)            not null,
    error_msg        varchar(2000),
    created_at       timestamp default now() not null,
    eval_started_at  timestamp,
    eval_finished_at timestamp
);

create index in_expressions_req_uid on expressions(req_uid);
alter table expressions add foreign key (user_id) REFERENCES users;
create index in_expressions_user_id on expressions(user_id);

create table agents
(
    name              varchar(255)  primary key,
    created_at        timestamp default now() not null,
    last_heartbeat_at timestamp default now() not null
);

create table operation_durations
(
    operation_name   varchar(255)  primary key,
    duration         numeric(7) not null,
    created_at       timestamp default now() not null,
    updated_at       timestamp default now() not null
);


create table sub_expressions
(
    id                 serial                  primary key,
    expression_id      integer,
    val1               numeric(50, 5),
    val2               numeric(50, 5),
    sub_expression_id1 integer,
    sub_expression_id2 integer,
    operation_name     varchar(255) not null,
    result             numeric(50, 5),
    agent_name         varchar(255),
    eval_started_at    timestamp,
    eval_finished_at   timestamp,
    is_last            boolean default false
);

alter table sub_expressions add foreign key (expression_id) REFERENCES expressions;

create index in_se_expression_id on sub_expressions(expression_id);
