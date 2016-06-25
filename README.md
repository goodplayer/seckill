# seckill
seckill experimental

# done

* [x] standard reduce
* [x] item query cache
* [x] hot item batch reduce
* [ ] inventory add back
* [ ] split inventory
* [ ] user qualifier
* [ ] Q&A filter
* [x] increase connection pool size
* [x] prepare statement for sql
* [x] use pg_try_advisory_xact_lock to reduce lock overhead. this may increase overall qps but not updating qps

# start postgresql

```
./initdb pghome
./initdb pghome2
```

```
sed -i 's/#port = 5432/port = 15432/g' pghome2/postgresql.conf
```

```
dstart -etoo -out ./pglog.log ./postgres -D ./pghome
dstart -etoo -out ./pglog2.log ./postgres -D ./pghome2
```

## db1

```
./psql -h localhost -p 5432 -d template1
 or
./createdb -h localhost -p 5432 $USERNAME
```

```
./psql -h localhost -p 5432

create user orderuser with password 'orderuser';
create database order_order owner orderuser;
GRANT ALL PRIVILEGES ON DATABASE order_order to orderuser;

create user inventoryuser with password 'invnetoryuser';
create database inventory owner inventoryuser;
GRANT ALL PRIVILEGES ON DATABASE inventory to inventoryuser;
```

##### order db init

```
./psql -U orderuser -d order_order -h 127.0.0.1 -p 5432
```

```
create table order_order (
    id bigint,
    user_id bigint,
    order_id bigint,
    create_time bigint,
    item_id bigint,
    status int,
    modify_time bigint,
    buy_quantity bigint,
    
    primary key(id)
);
CREATE UNIQUE INDEX order_order_order_id ON order_order (order_id);
```

##### inventory db init

```
./psql -U inventoryuser -d inventory -h 127.0.0.1 -p 5432
```

```
create table item_inventory (
    id bigint,
    item_id bigint,
    quantity bigint,
    status int,
    create_time bigint,
    modify_time bigint,
    parent_id bigint,
    root_id bigint,
    user_id bigint,

    primary key(id)
);
CREATE INDEX item_inventory_item_id_index ON item_inventory (item_id);
```

## db2

```
./psql -h localhost -p 15432 -d template1
 or
./createdb -h localhost -p 15432 $USERNAME
```

```
./psql -h localhost -p 15432

create user inventoryuser2 with password 'invnetoryuser2';
create database inventory2 owner inventoryuser2;
GRANT ALL PRIVILEGES ON DATABASE inventory2 to inventoryuser2;
```

##### inventory db2 init

```
./psql -U inventoryuser2 -d inventory2 -h 127.0.0.1 -p 15432
```

```
create table item_inventory (
    id bigint,
    item_id bigint,
    quantity bigint,
    status int,
    create_time bigint,
    modify_time bigint,
    parent_id bigint,
    root_id bigint,
    user_id bigint,

    primary key(id)
);
CREATE INDEX item_inventory_item_id_index ON item_inventory (item_id);
```
