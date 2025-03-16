create extension if not exists "uuid-ossp";

alter table schema_name.orders
    alter column id set default uuid_generate_v4();
