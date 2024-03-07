alter database dev set timezone to 'Europe/Moscow';

create table if not exists users(
    id bigserial primary key,
    email text unique not null,
    password text not null,
    first_name text not null,
    last_name text not null default '',
    last_pass_reset timestamp default current_timestamp,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);