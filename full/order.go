package full

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

func generateOrderId(userId int64) (int64, error) {
	return rand.Int63(), nil
}

func SaveOrder(orderId, userId, itemId, quantity int64) error {
	tag, err := order_pg_pool.Exec(
		"insert into order_order(id, user_id, order_id, create_time, item_id, status, modify_time, buy_quantity) values($1, $2, $3, $4, $5, $6, $7, $8)",
		orderId, userId, orderId, time.Now().Unix(), itemId, 0, time.Now().Unix(), quantity,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("insert order row affected is not 0. result is: " + strconv.Itoa(int(tag.RowsAffected())))
	}
	return nil
}
