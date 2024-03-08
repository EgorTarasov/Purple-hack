alter database dev set timezone to 'Europe/Moscow';

create table if not exists session(
    id uuid primary key,
    prompt_ids bigint[] default array[]::bigint[],
    created_at timestamp default current_timestamp
);

create table if not exists prompt(
    id bigserial primary key,
    fk_session_id uuid references session(id),
    body text not null,
    user_agent text not null,
    created_at timestamp default current_timestamp
);

create table if not exists response(
    id bigserial primary key,
    session_id uuid references session(id),
    prompt_id bigint references prompt(id),
    body text not null,
    created_at timestamp default current_timestamp
);