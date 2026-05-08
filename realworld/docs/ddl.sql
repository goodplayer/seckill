create table seckill_user
(
    user_id      uuid primary key not null,
    username     varchar(256)     not null,
    time_created int8             not null,
    time_updated int8             not null
);

insert into seckill_user(user_id, username, time_created, time_updated)
values ('408f6222-c435-4535-9de8-bf4ca22a79bc', 'demo_buyer', 1778208790000, 1778208790000),
       ('772a035b-3ca1-49b9-a047-ab12fc04b2c4', 'demo_seller', 1778208790000, 1778208790000);

create table seckill_order
(
    order_id     uuid primary key not null,
    user_id      uuid             not null,
    seller_id    uuid             not null,
    order_item   uuid             not null,
    amount       int8             not null,
    unit_price   int8             not null, -- same as unit_price in item table
    total_price  int8             not null, -- value = amount * total_price
    order_status int8             not null, -- -1 - invisible, 0 - invalid status, 1 - created, 2 - paid, ...
    time_created int8             not null,
    time_updated int8             not null
);

CREATE INDEX idx_order_user_id ON seckill_order (user_id);

create table seckill_item
(
    item_id      uuid primary key not null,
    seller_id    uuid             not null,
    item_name    varchar(256)     not null,
    unit_price   int8             not null, -- the value here = 100 * the original price which has 2-digits
    description  text             not null,
    status       int8             not null, -- 0 is normal
    inventory_id uuid             not null,
    time_created int8             not null,
    time_updated int8             not null
);

CREATE INDEX idx_item_seller_id ON seckill_item (seller_id);

insert into seckill_item(item_id, seller_id, item_name, unit_price, description, status, inventory_id, time_created,
                         time_updated)
VALUES ('f7e61eef-6f7e-4f50-9f75-19c64b06a7f2', '772a035b-3ca1-49b9-a047-ab12fc04b2c4', 'demo_item', '3',
        'This is a demo item.', 0, '86f0127a-ab78-45e6-acb4-26fd6a7f5f07', 1778208790000, 1778208790000);

create table seckill_inventory
(
    inventory_id      uuid primary key not null,
    total_stock       int8             not null,
    withholding_stock int8             not null,
    time_created      int8             not null,
    time_operated     int8             not null
);

create table seckill_inventory_order
(
    inventory_order_record_id uuid primary key not null,
    inventory_id              uuid             not null,
    order_id                  uuid             not null,
    amount                    int8             not null,
    status                    int8             not null, -- 0 - withheld, 1 - deducted, 2 - released
    time_created              int8             not null,
    time_updated              int8             not null,
    UNIQUE (order_id)
);

insert into seckill_inventory (inventory_id, total_stock, withholding_stock, time_created, time_operated)
values ('86f0127a-ab78-45e6-acb4-26fd6a7f5f07', 11111111111111111, 0, 1778208790000, 1778208790000);
