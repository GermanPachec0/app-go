package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GermanPachec0/app-go/domain"
	"github.com/GermanPachec0/app-go/repository"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	listenAddr string
	userRepo   domain.UserRepository
}

func NewAPI(ctx context.Context, pool *pgxpool.Pool) *APIServer {
	userRepo := repository.NewPostgresUser(pool)

	return &APIServer{
		listenAddr: "8080",
		userRepo:   userRepo,
	}
}

func (a *APIServer) Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/v1/users/{uid}", a.getUserByIdHandler).Methods("GET")
	return r
}

func (a *APIServer) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}
