package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/jackc/pgx.v2"

	"import.moetang.info/go/lib/gin-startup"

	"github.com/goodplayer/seckill/full"
)

func main() {
	full.Init(&full.Config{
		InventoryPgConfig: pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Port:     5432,
				Database: "inventory",
				User:     "inventoryuser",
				Password: "inventoryuser",
			},
			MaxConnections: 200,
		},
		OrderPgConfig: pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Port:     5432,
				Database: "order_order",
				User:     "orderuser",
				Password: "orderuser",
			},
			MaxConnections: 200,
		},
	})

	g := gin_startup.NewGinStartup()
	g.Custom(func(r *gin.Engine) {
		r.GET("/reducing", func(c *gin.Context) {

			itemIdStr, ok := c.GetQuery("item_id")
			userIdStr, ok := c.GetQuery("user_id")

			itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
			if !ok || err != nil {
				c.JSON(http.StatusBadRequest, `{"error_msg":"itemId format error"}`)
				return
			}
			userId, err := strconv.ParseInt(userIdStr, 10, 64)
			if !ok || err != nil {
				c.JSON(http.StatusBadRequest, `{"error_msg":"userId format error"}`)
				return
			}

			req := &full.CreateOrderReq{
				UserId:      userId,
				ItemId:      itemId,
				BuyQuantity: 1,
			}
			resp, err := full.CreateOrder(req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			} else {
				c.JSON(http.StatusInternalServerError, `{"order_id":`+strconv.FormatInt(resp.OrderId, 10)+`}`)
				return
			}
		})
	})
	g.EnableHttp("tcp://127.0.0.1:7649")
	g.Start()

	cccc := make(chan bool)
	<-cccc
}
