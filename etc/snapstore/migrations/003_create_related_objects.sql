-- Write your migrate up statements here
create table related_objects(
    id text not null,
    type text not null,
    data jsonb not null,
    date_created timestamptz not null default now(),
    date_updated timestamptz
);

create unique index related_objects_primary_idx on related_objects(id,type)
create index related_objects_id_idx on related_objects(id);
create index related_objects_type_idx on related_objects(type);
create index related_objects_date_created_idx on related_objects(date_created);
create index related_objects_date_updated_idx on related_objects(date_updated);

---- create above / drop below ----

drop table related_objects cascade;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
