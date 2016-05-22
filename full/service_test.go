package full_test

import (
	"math/rand"
	"testing"
	"time"

	"gopkg.in/jackc/pgx.v2"

	"github.com/goodplayer/seckill/full"
	"github.com/goodplayer/seckill/global"
)

func init() {
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
}

func TestCreateOrder(t *testing.T) {
	req := &full.CreateOrderReq{
		ItemId:      rand.Int63n(global.TOTAL_ITEM_ID_COUNT) + global.START_ITEM_ID,
		UserId:      1,
		BuyQuantity: 1,
	}
	resp, err := full.CreateOrder(req)
	if err != nil {
		t.Fatal("create error fail.", err)
	}
	t.Log("create order id:", resp.OrderId)
}

func BenchmarkCreateOrder(b *testing.B) {
	req := &full.CreateOrderReq{
		ItemId:      rand.Int63n(global.TOTAL_ITEM_ID_COUNT) + global.START_ITEM_ID,
		UserId:      1,
		BuyQuantity: 1,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req.ItemId = rand.Int63n(global.TOTAL_ITEM_ID_COUNT) + global.START_ITEM_ID
		_, err := full.CreateOrder(req)
		if err != nil {
			b.Fatal("create error fail.", err)
		}
	}
}

func BenchmarkParallelCreateOrder(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := &full.CreateOrderReq{
				ItemId:      rand.Int63n(global.TOTAL_ITEM_ID_COUNT) + global.START_ITEM_ID,
				UserId:      1,
				BuyQuantity: 1,
			}
			_, err := full.CreateOrder(req)
			if err != nil {
				b.Fatal("create error fail.", err)
			}
		}
	})
}

func TestCreateOrderNTimes(t *testing.T) {
	var N = 1000

	var startTime = time.Now().UnixNano()

	for i := 0; i < N; i++ {
		req := &full.CreateOrderReq{
			ItemId:      rand.Int63n(global.TOTAL_ITEM_ID_COUNT) + global.START_ITEM_ID,
			UserId:      1,
			BuyQuantity: 1,
		}
		_, err := full.CreateOrder(req)
		if err != nil {
			t.Fatal("create error fail.", err)
		}
	}

	var endTime = time.Now().UnixNano()

	var during = endTime - startTime

	t.Log("qps:", float64(N)/float64(during)*1000000000)
}
