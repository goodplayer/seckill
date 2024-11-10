# Seckill PoC

业务场景

* 库存扣减
* 优惠券抢购
* 红包抢购

业务特点

* 单行更新
* 流量大
* 准确性要求强

PoC业务目标

* 1000w users
* 1000w inventory items
* 100w tps total / 10w tps success

参考数值

* 某业务180w tps/ 30w tps, 深度 2 billion

## 测试环境和要求

环境配置

```text
xeon e-2286M 4c/8g
postgresql 17
Micro 7450 pro 960g - 数据盘
Intel Optane M10 16g - wal盘
```

pg参数配置

```text
max_connections = 300
tcp_keepalives_idle = 300
tcp_keepalives_interval = 3
tcp_keepalives_count = 5

shared_buffers = 2GB
huge_pages = try
temp_buffers = 32MB
work_mem = 3495kB
maintenance_work_mem = 512MB
effective_io_concurrency = 200
max_worker_processes = 4
max_parallel_workers_per_gather = 2
max_parallel_maintenance_workers = 2
max_parallel_workers = 4

wal_level = logical
wal_buffers = 16MB
checkpoint_timeout = 5min
checkpoint_completion_target = 0.9
min_wal_size = 2GB
max_wal_size = 8GB
summarize_wal = on

max_wal_senders = 16
max_replication_slots = 16

random_page_cost = 1.1
effective_cache_size = 6GB
default_statistics_target = 100

idle_in_transaction_session_timeout = 60min

```

DDL

```text
-- inventory item definition used as core model
create table inventory_item (
	item_id int8 not null,
	item_2nd_id int8 not null,
	item_3rd_id int8 not null,
	inventory_id int8 not null,
	time_created int8 not null,
	time_updated int8 not null,
	primary key (item_id, item_2nd_id, item_3rd_id)
);

create unique index on inventory_item(inventory_id);

-- operation log on inventory changes
create table inventory_operation (
	item_id int8 not null,
	inventory_id int8 not null,
	operation varchar not null,
	quantity int8 not null,
	time_created int8 not null,
	primary key(item_id)
);

-- inventory quantity used as core quantity management
create table inventory_quantity (
	inventory_id int8 not null primary key,
	partition int8 not null,
	quantity int8 not null,
	quantity_2nd int8 not null,
	time_created int8 not null
);

create unique index on inventory_quantity (inventory_id, partition);

-- inventory quantity usage used as tracking inventory quantity usage
create table inventory_quantity_usage (
	inventory_id int8 not null,
	usage_id varchar not null,
	usage_type varchar not null,
	quantity int8 not null,
	time_created int8 not null,
	primary key (inventory_id, usage_id, usage_type)
);

```

业务行为

```text
1. 准备的库存数据存放在inventory_quantity表中
2. 扣减操作需要在更新inventory_quantity表的同时，插入扣减记录到inventory_quantity_usage表
3. inventory_item和inventory_operation表仅作为业务描述使用。测试中不需要操作。
```

数据记录

```text
吞吐量
    总吞吐
    真实吞吐
数据大小（使用psql，命令\l+）
    测试前
    测试后
系统数据
    大约cpu使用
    大约磁盘bandwidth、iops、util
```

测试方式及注意

```text
1. 使用pgbench进行测试
2. 每次完成一轮测试，进行vacuum，并等待checkpoint完成
```

## 测试场景和结果

### 0. 准备数据

```text
-- 100 million inventory items
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(1,10000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(10000001,20000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(20000001,30000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(30000001,40000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(40000001,50000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(50000001,60000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(60000001,70000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(70000001,80000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(80000001,90000000), 0, 100000000, 100000000, 1731207910000;
insert into inventory_quantity (inventory_id, partition, quantity, quantity_2nd, time_created) select generate_series(90000001,100000000), 0, 100000000, 100000000, 1731207910000;
```

清理数据

```text
update inventory_quantity set quantity = 100000000 where quantity < 100000000;
```

### 1. 单行直接更新

脚本`nano test.sql`

```text
\set q random(1,10)
BEGIN;
WITH updated_rows AS (
update inventory_quantity set quantity = quantity - :q where inventory_id = 30000000 and quantity >= :q returning inventory_id
) select EXISTS(select inventory_id from updated_rows) as has_updated \gset
\if :has_updated
insert into inventory_quantity_usage (inventory_id, usage_id, usage_type, quantity, time_created) values (30000000, random(1,100000000000), 'D', :q, 1731207910000);
\endif
END;
```

结果获取

```text
./pgbench -U admin -d inventory -M prepared -n -r -P 1 -f ./test.sql -c 64 -j 64 -T 100
iostat -x 4
select count(*) from inventory_quantity_usage ;
```

清理脚本

```text
truncate inventory_quantity_usage ;
update inventory_quantity set quantity = 100000000 where quantity < 100000000;
vacuum full ;
```

测试结果

```text
tps = 3005.225247 (without initial connection time)
orders = 300548
3005.48tps

iostat wal write:
3000iops
48MB/s

data size = 11G
```

### 2. 单行 + advisory

脚本`nano test.sql`

```text
\set q random(1,10)
BEGIN;
WITH updated_rows AS (
update inventory_quantity set quantity = quantity - :q where inventory_id = 30000000 and quantity >= :q and pg_try_advisory_xact_lock(30000000) returning inventory_id
) select EXISTS(select inventory_id from updated_rows) as has_updated \gset
\if :has_updated
insert into inventory_quantity_usage (inventory_id, usage_id, usage_type, quantity, time_created) values (30000000, random(1,100000000000), 'D', :q, 1731207910000);
\endif
END;
```

结果获取

```text
./pgbench -U admin -d inventory -M prepared -n -r -P 1 -f ./test.sql -c 64 -j 64 -T 100
iostat -x 4
select count(*) from inventory_quantity_usage ;
```

清理脚本

```text
truncate inventory_quantity_usage ;
update inventory_quantity set quantity = 100000000 where quantity < 100000000;
vacuum full ;
```

测试结果

```text
tps = 48653.240665 (without initial connection time)
orders = 24420
244.2tps

iostat wal write:
280iops
4.2MB/s

data size = 11G
```

### 3. 多行直接更新

脚本`nano test.sql`

```text
\set q random(1,10)
\set id random(1,100000000)
BEGIN;
WITH updated_rows AS (
update inventory_quantity set quantity = quantity - :q where inventory_id = :id and quantity >= :q returning inventory_id
) select EXISTS(select inventory_id from updated_rows) as has_updated \gset
\if :has_updated
insert into inventory_quantity_usage (inventory_id, usage_id, usage_type, quantity, time_created) values (:id, random(1,100000000000), 'D', :q, 1731207910000);
\endif
END;
```

结果获取

```text
./pgbench -U admin -d inventory -M prepared -n -r -P 1 -f ./test.sql -c 64 -j 64 -T 100
iostat -x 4
select count(*) from inventory_quantity_usage ;
```

清理脚本

```text
truncate inventory_quantity_usage ;
update inventory_quantity set quantity = 100000000 where quantity < 100000000;
vacuum full ;
```

测试结果

```text
tps = 4949.788535 (without initial connection time)
orders = 495340
4953.4tps

iostat wal write:
220iops
140MB/s

data size = 11G
```

### 4. 多行 + advisory

TBD

### 5. 单行拆分多行直接更新

TBD

### 6. 单行拆分多行 + advisory

TBD

### 7. 多行直接更新 + 顺序调整

TBD

## References

* [聊一聊双十一背后的技术 - 不一样的秒杀技术, 裸秒](https://github.com/digoal/blog/blob/master/201611/20161117_01.md)
* [HTAP数据库 PostgreSQL 场景与性能测试之 28 - (OLTP) 高并发点更新](https://github.com/digoal/blog/blob/master/201711/20171107_29.md)
* [HTAP数据库 PostgreSQL 场景与性能测试之 30 - (OLTP) 秒杀 - 高并发单点更新](https://github.com/digoal/blog/blob/master/201711/20171107_31.md)

