package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brix101/linkr/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type api struct {
	DB db.Queries
}

func New(ctx context.Context, pool *pgxpool.Pool) *api {
	q := db.New(pool)

	return &api{
		DB: *q,
	}
}

func (a *api) Server(port string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: a.loadRoutes(),
	}
}

func (a *api) loadRoutes() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/{code}", a.getLink)
	r.HandleFunc("POST /link", a.createLink)

	return r
}
