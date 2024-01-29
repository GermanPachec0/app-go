package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/GermanPachec0/app-go/domain"
	jwt "github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Email string `json:"email"`
	Type  string `json:"type"`
	Token string `json:"token"`
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	logReq := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
		s.errorResponse(w, r, 500, err)
		return
	}
	user, err := s.userRepo.GetByEmail(ctx, logReq.Email)
	if err != nil {
		s.errorResponse(w, r, 404, err)
		return
	}
	if !user.ValidatePassword(logReq.Password) {
		s.errorResponse(w, r, 401, err)
		return
	}

	token, err := createJWT(&user)
	if err != nil {
		s.errorResponse(w, r, 500, err)
	}

	res := LoginResponse{
		Email: user.Email,
		Type:  user.Type,
		Token: token,
	}

	WriteJson(w, 200, res)
}

func createJWT(user *domain.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":   15000,
		"memeber_uid": user.Uuid,
		"role":        user.Type,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}
