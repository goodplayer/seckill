package full_test

import (
	"testing"

	"gopkg.in/jackc/pgx.v2"

	"github.com/goodplayer/seckill/full"
)

func TestSaveOrder(t *testing.T) {
	full.Init(&full.Config{
		OrderPgConfig: pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Port:     5432,
				Database: "order_order",
				User:     "orderuser",
				Password: "orderuser",
			},
			MaxConnections: 20,
		},
	})

	t.Log(full.SaveOrder(1, 1, 1))
}
