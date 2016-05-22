package full

import (
	"errors"
)

func QueryInventory(itemId int64) (int64, error) {
	row, err := inventory_pg_pool.Query("select quantity from item_inventory where item_id = $1 and status = 0", itemId)
	if err != nil {
		return 0, err
	}
	defer row.Close()
	if row.Next() {
		var quantity int64
		err = row.Scan(&quantity)
		if err != nil {
			return 0, err
		}
		return quantity, nil
	} else {
		return 0, errors.New("item not found or item status abnormal.")
	}
}

func ReduceInventory(itemId, quantity int64) (int64, error) {
	r, err := inventory_pg_pool.Query("update item_inventory set quantity = quantity - $1 where item_id = $2 and status = 0 and quantity >= $3 returning quantity", quantity, itemId, quantity)
	if err != nil {
		return 0, err
	} else {
		defer r.Close()
		if r.Next() {
			var newQuantity int64
			err = r.Scan(&newQuantity)
			if err != nil {
				return 0, err
			}
			return newQuantity, nil
		} else {
			return 0, errors.New("no item inventory reduced.")
		}
	}
}

func UpdateQuantityCache(itemId, quantity int64) error {
	//TODO quantity cache
	return nil
}
