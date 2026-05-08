-- add cleanup scripts
truncate seckill_inventory_order;
update seckill_inventory set total_stock = 11111111111111111, withholding_stock = 0;
truncate seckill_order;
VACUUM FULL;
