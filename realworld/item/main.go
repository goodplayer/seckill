package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/goodplayer/seckill/realworld/item/item"
	"github.com/goodplayer/seckill/realworld/shared"
)

func main() {
	cfg, err := pgxpool.ParseConfig(shared.ItemDatabaseConnectionString)
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

	r.Get("/item/{item_id}", func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "item_id")
		if itemId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		i, err := item.QueryItemById(itemId, pool)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if data, err := json.Marshal(i); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)
		}
	})

	if err := http.ListenAndServe(":3002", r); err != nil {
		panic(err)
	}
}
