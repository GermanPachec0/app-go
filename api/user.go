package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *APIServer) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	uidStr := mux.Vars(r)["uid"]
	if uidStr == "" {
		a.errorResponse(w, r, 500, errors.New("no params found"))
		return
	}
	uid, err := uuid.ParseBytes([]byte(uidStr))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	user, err := a.userRepo.GetById(ctx, uid)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	WriteJson(w, 200, user)
}
