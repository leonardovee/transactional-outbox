create type order_status as enum (
    'created',
    'approved',
    'ready',
    'shipped',
    'arrived'
);

create table orders (
    id varchar(36) primary key,
    aggregate_id varchar(36) not null,
    status order_status,
    total int not null,
    created_at timestamp default now()
);

create index orders_aggregate_id on orders(aggregate_id);
