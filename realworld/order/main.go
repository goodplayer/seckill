package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"resty.dev/v3"

	"github.com/goodplayer/seckill/realworld/order/order"
	"github.com/goodplayer/seckill/realworld/shared"
)

func main() {
	cfg, err := pgxpool.ParseConfig(shared.OrderDatabaseConnectionString)
	if err != nil {
		panic(err)
	}
	cfg.MaxConns = 50
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//r.Use(middleware.Logger) // prevent mass log output that impact performance
	r.Use(middleware.Recoverer)

	orderService := &order.OrderService{
		DependencyService: &order.RemoteDependencyService{
			Client: resty.NewWithTransportSettings(&resty.TransportSettings{
				DialerTimeout:       10 * time.Second,
				DialerKeepAlive:     10 * time.Second,
				IdleConnTimeout:     90 * time.Second,
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
			}),
			Pool: pool,
		},
	}

	r.Get("/echo", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("welcome"))
	})

	r.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		req := new(shared.PlaceOrderRequest)
		if err := ParseJsonBody(req, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if id, err := orderService.PlaceOrder(req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"order_id": id,
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})

	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}

func ParseJsonBody(obj any, req *http.Request) error {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}
