-- Write your migrate up statements here

alter table response rename column query_id to fk_query_id;

---- create above / drop below ----

alter table response rename column fk_query_id to query_id;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
