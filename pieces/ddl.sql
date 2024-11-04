create table inventory
(
    item_id      int8 not null,
    quantity     int8 not null,
    parent_id    int8 not null,
    time_created int8 not null,
    time_updated int8 not null,
    CONSTRAINT inventory_pk PRIMARY KEY (item_id)
);
