create table if not exists users
(
    id       varchar(36) primary key,
    login    varchar(50) unique not null,
    password varchar(255)       not null,
    created  timestamp          not null
);

create table if not exists documents
(
    id         varchar(36) primary key,
    name       varchar(255) not null,
    mime       varchar(100),
    file       boolean      not null,
    public     boolean      not null,
    created    timestamp    not null,
    grant_list text[],
    owner      varchar(36)  not null,
    data       bytea,
    json       text,
    foreign key (owner) references users (id)
);

create index if not exists idx_documents_owner on documents (owner);
create index if not exists idx_documents_created on documents (created);