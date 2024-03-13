alter database cbr set timezone to 'Europe/Moscow';

create table if not exists session(
    id uuid primary key,
    created_at timestamp default current_timestamp
);

create table if not exists query(
    id bigserial primary key,
    fk_session_id uuid references session(id),
    body text not null,
    model text not null,
    user_agent text not null,
    created_at timestamp default current_timestamp
);

create table if not exists response(
    id bigserial primary key,
    fk_session_id uuid references session(id),
    query_id bigint references query(id),
    body text not null,
    context jsonb not null,
    created_at timestamp default current_timestamp
);

create table if not exists "user"(
    id bigserial primary key,
    email text unique not null,
    password text not null,
    created_at timestamp default current_timestamp
);

create table if not exists users_sessions(
    fk_user_id bigint references "user"(id),
    fk_session_id uuid references session(id),
    created_at timestamp default current_timestamp,
    primary key (fk_user_id, fk_session_id)
);
