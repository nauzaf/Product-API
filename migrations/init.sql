create table products (
    id serial primary key,
    name varchar(255) NOT NULL,
    inventory int NOT NULL DEFAULT 0,
    expiry timestamp,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
)