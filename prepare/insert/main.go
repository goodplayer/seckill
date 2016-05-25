package main

import (
	"log"
	"math/rand"
	"time"

	"gopkg.in/jackc/pgx.v2"
)

func main() {
	config := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "127.0.0.1",
			Port:     5432,
			Database: "inventory",
			User:     "inventoryuser",
			Password: "inventoryuser",
		},
		MaxConnections: 20,
	}

	pool, err := pgx.NewConnPool(config)

	if err != nil {
		log.Fatalln("new conn pool error.", err)
	}

	rand.Seed(time.Now().UnixNano())

	var itemInventory *ItemInventory = &ItemInventory{}

	var tx *pgx.Tx
	var e error
	for i := 0; i < 1000; i++ {
		tx, e = pool.Begin()
		if e != nil {
			log.Fatalln("begin tx error.", e)
		}
		for k := 0; k < 10000; k++ {
			itemInventory.ItemId = int64(i*10000+k) + 10000000000
			itemInventory.Quantity = 999999999999
			itemInventory.Status = 0
			create(tx, itemInventory)
		}
		e = tx.Commit()
		if e != nil {
			log.Fatalln("commit error.", e)
		}
		log.Println((i + 1) * 10000)
	}

	tx, e = pool.Begin()
	if e != nil {
		log.Fatalln("begin tx error for itemId=2000000000", e)
	}
	itemInventory.ItemId = 2000000000
	itemInventory.Quantity = 10
	itemInventory.Status = 0
	create(tx, itemInventory)
	e = tx.Commit()
	if e != nil {
		log.Fatalln("commit error of itemId=2000000000.", e)
	}

	log.Println("done! 10000000")
}

func create(tx *pgx.Tx, itemInventory *ItemInventory) {
	sqlString := "insert into item_inventory(id, item_id, quantity, status, create_time, modify_time, parent_id, root_id, user_id) values($1, $2, $3, $4, $5, $6, $7, $8, $9);"
	_, e := tx.Exec(sqlString, itemInventory.ItemId, itemInventory.ItemId, itemInventory.Quantity, itemInventory.Status, time.Now().Unix(), time.Now().Unix(), 0, 0, rand.Int63())
	if e != nil {
		log.Fatalln("insert error.", e)
	}
}

type ItemInventory struct {
	ItemId   int64
	Quantity int64
	Status   int
}
