package main

import (
	"log"

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

	config2 := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "127.0.0.1",
			Port:     5432,
			Database: "inventory2",
			User:     "inventoryuser2",
			Password: "inventoryuser2",
		},
		MaxConnections: 20,
	}

	pool, err := pgx.NewConnPool(config)
	if err != nil {
		log.Fatalln("new conn pool error.", err)
	}
	pool2, err := pgx.NewConnPool(config2)
	if err != nil {
		log.Fatalln("new conn pool error.", err)
	}

	_, err = pool.Exec("update item_inventory set quantity = 999999999999 where quantity < 999999999999;")
	if err != nil {
		log.Fatalln("update item_inventory error.", err)
	}

	_, err = pool.Exec("update item_inventory set quantity = 10 where item_id = 2000000000;")
	if err != nil {
		log.Fatalln("update item_inventory error.", err)
	}

	_, err = pool.Exec("update item_inventory set quantity = 888888888 where item_id = 3000000000;")
	if err != nil {
		log.Fatalln("update item_inventory error.", err)
	}

	_, err = pool2.Exec("update item_inventory set quantity = 888888888 where item_id = 3000000000;")
	if err != nil {
		log.Fatalln("update item_inventory error.", err)
	}
}
