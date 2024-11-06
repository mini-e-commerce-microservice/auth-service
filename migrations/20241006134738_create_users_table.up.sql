CREATE TABLE users
(
    id                bigserial primary key,
    email             varchar(255) not null,
    password          varchar(255) not null,
    is_email_verified boolean      not null default false,
    created_at        timestamp    not null,
    updated_at        timestamp    not null,
    deleted_at        timestamp,
    trace_parent       varchar
);