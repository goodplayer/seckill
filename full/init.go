package full

import (
	"math/rand"
	"time"

	"gopkg.in/jackc/pgx.v2"
)

var (
	order_pg_pool     *pgx.ConnPool
	inventory_pg_pool *pgx.ConnPool
)

type Config struct {
	OrderPgConfig      pgx.ConnPoolConfig
	InventoryPgConfig  pgx.ConnPoolConfig
	InventoryPg2Config pgx.ConnPoolConfig
}

func Init(config *Config) error {
	rand.Seed(time.Now().UnixNano()) // for generate order id
	var err error
	order_pg_pool, err = pgx.NewConnPool(config.OrderPgConfig)
	if err != nil {
		return err
	}
	inventory_pg_pool, err = pgx.NewConnPool(config.InventoryPgConfig)
	if err != nil {
		return err
	}
	return nil
}

// for hot item to be cached
func PreLoadInventoryData() {
	//TODO
}
