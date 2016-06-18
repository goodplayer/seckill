package full

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	orderSaveKey string
)

func prepareOrderSql() {
	orderSave := "insert into order_order(id, user_id, order_id, create_time, item_id, status, modify_time, buy_quantity) values($1, $2, $3, $4, $5, $6, $7, $8)"
	orderSaveMd5 := md5.Sum([]byte(orderSave))
	orderSaveKey = "orderSave_" + hex.EncodeToString(orderSaveMd5[:])
	_, err := order_pg_pool.Prepare(orderSaveKey, orderSave)
	if err != nil {
		log.Fatalln("prepare sql order - orderSave error", err)
	} else {
		log.Println("prepare sql order - orderSave - key:", orderSaveKey, "sql:", orderSave)
	}
}

func generateOrderId(userId int64) (int64, error) {
	return rand.Int63(), nil
}

func SaveOrder(orderId, userId, itemId, quantity int64) error {
	tag, err := order_pg_pool.Exec(
		orderSaveKey,
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
