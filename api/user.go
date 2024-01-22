package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *APIServer) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	idStr := mux.Vars(r)["id"]
	uid, err := uuid.FromBytes([]byte(idStr))
	if err != nil {
		a.errorResponse(w, r, 500, err)
	}
	user, err := a.userRepo.GetById(ctx, uid)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	WriteJson(w, 200, user)
}
