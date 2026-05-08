package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/goodplayer/seckill/realworld/inventory/inventory"
	"github.com/goodplayer/seckill/realworld/shared"
)

func main() {
	cfg, err := pgxpool.ParseConfig(shared.InventoryDatabaseConnectionString)
	if err != nil {
		panic(err)
	}
	cfg.MaxConns = 50
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//r.Use(middleware.Logger) // prevent mass log output that impact performance
	r.Use(middleware.Recoverer)

	r.Post("/inventory/withhold", func(w http.ResponseWriter, r *http.Request) {
		req := new(shared.InventoryWithholdRequest)
		if err := ParseJsonBody(req, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := inventory.WithholdInventory(req, pool); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(":3003", r); err != nil {
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
