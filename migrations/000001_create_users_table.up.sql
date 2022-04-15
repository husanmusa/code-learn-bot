create table if not exists users (
    id bigserial not null primary key,
    chat_id bigint not null,
    input_name varchar(32),
    is_banned bool default false
)