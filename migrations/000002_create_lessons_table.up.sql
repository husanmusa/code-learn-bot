create table if not exists lessons
(
    id               bigserial not null primary key,
    chat_id          bigint    not null,
    message_id       bigint    not null,
    number_of_lesson int       not null,
    number_of_part   int       not null,
    type_of_part     varchar(24)
);