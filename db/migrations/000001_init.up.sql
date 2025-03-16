create schema if not exists schema_name;

create table if not exists schema_name.orders
(
    id       uuid              not null
        constraint orders_pk
            primary key,
    item     text              not null,
    quantity integer default 0 not null
);