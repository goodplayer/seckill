package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/goodplayer/seckill/full"
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
			MaxConnections: 10,
		},
		InventoryPg2Config: pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Port:     15432,
				Database: "inventory2",
				User:     "inventoryuser2",
				Password: "inventoryuser2",
			},
			MaxConnections: 10,
		},
		OrderPgConfig: pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Port:     5432,
				Database: "order_order",
				User:     "orderuser",
				Password: "orderuser",
			},
			MaxConnections: 10,
		},
	})
}

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	g := gin.New()
	{
		//r.Use(gin.Logger(), gin.Recovery())
		g.GET("/reducing", func(c *gin.Context) {

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
				c.JSON(http.StatusOK, `{"order_id":`+strconv.FormatInt(resp.OrderId, 10)+`}`)
				return
			}
		})
	}
	go func() {
		u, err := url.Parse("tcp://127.0.0.1:7649")
		if err != nil {
			panic(err)
		}
		if gin.IsDebugging() {
			log.Printf("[GIN-debug] Listening and serving HTTP on %s\n", u.Host)
		}
		defer func() {
			if err != nil && gin.IsDebugging() {
				log.Printf("[GIN-debug] [ERROR] %v\n", err)
			}
		}()

		server := &http.Server{
			Addr:    u.Host,
			Handler: g,
		}

		err = server.ListenAndServe()
	}()

	cccc := make(chan bool)
	<-cccc
}
