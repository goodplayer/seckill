select count(*) as order_count_total from seckill_order;
select count(*) as order_count_activated from seckill_order where order_status = 1;
select withholding_stock from seckill_inventory where inventory_id = '86f0127a-ab78-45e6-acb4-26fd6a7f5f07';
