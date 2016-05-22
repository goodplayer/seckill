package full_test

import (
	"testing"

	"gopkg.in/jackc/pgx.v2"

	"github.com/goodplayer/seckill/full"
)

func TestQueryInventory(t *testing.T) {
	full.Init(&full.Config{
		InventoryPgConfig: pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Port:     5432,
				Database: "inventory",
				User:     "inventoryuser",
				Password: "inventoryuser",
			},
			MaxConnections: 20,
		},
	})

	t.Log(full.QueryInventory(10000000001))
}

func TestReduceInventory(t *testing.T) {
	full.Init(&full.Config{
		InventoryPgConfig: pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Port:     5432,
				Database: "inventory",
				User:     "inventoryuser",
				Password: "inventoryuser",
			},
			MaxConnections: 20,
		},
	})

	t.Log(full.ReduceInventory(10000000001, 1))
}
