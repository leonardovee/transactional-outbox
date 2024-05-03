create table outbox (
    id varchar(36) primary key,
    aggregate_id varchar(36) not null,
    aggregate_type varchar(255) not null,
    type varchar(255) not null,
    payload jsonb not null,
    created_at timestamp default now()
);

create index outbox_aggregate_id on outbox(aggregate_id);
