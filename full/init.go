package full

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	order_pg_pool      *pgxpool.Pool
	inventory_pg_pool  *pgxpool.Pool
	inventory_pg2_pool *pgxpool.Pool
)

type Config struct {
	OrderPgConfig      pgxpool.Config
	InventoryPgConfig  pgxpool.Config
	InventoryPg2Config pgxpool.Config
}

func Init(config *Config) error {
	rand.Seed(time.Now().UnixNano()) // for generate order id
	var err error
	order_pg_pool, err = pgxpool.NewWithConfig(context.Background(), &config.OrderPgConfig)
	if err != nil {
		return err
	}
	inventory_pg_pool, err = pgxpool.NewWithConfig(context.Background(), &config.InventoryPgConfig)
	if err != nil {
		return err
	}
	inventory_pg2_pool, err = pgxpool.NewWithConfig(context.Background(), &config.InventoryPg2Config)
	if err != nil {
		return err
	}
	prepareInventorySql()
	prepareOrderSql()
	return nil
}

// for hot item to be cached
func PreLoadInventoryData() {
	//TODO
}
