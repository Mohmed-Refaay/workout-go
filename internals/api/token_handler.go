package api

import (
	"backend-go/internals/store"
	"backend-go/internals/tokens"
	"backend-go/internals/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TokenHandler struct {
	userStore  store.UserStore
	tokenStore store.TokenStore
	logger     *log.Logger
}

type authenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type getUserFromTokenRequest struct {
	Token string `json:"token"`
	Scope string `json:"scope"`
}

func NewTokenHandler(userStore store.UserStore, tokenStore store.TokenStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		userStore:  userStore,
		tokenStore: tokenStore,
		logger:     logger,
	}
}

func (th *TokenHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	payload := &authenticateRequest{}

	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		th.logger.Printf("Error: Authentication decoding %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error!"})
		return
	}

	user, err := th.userStore.GetUserByUsername(payload.Username)
	if err != nil {
		th.logger.Printf("Error: Authentication GetUserByUsername %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error!"})
		return
	}
	if user == nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid credentials!"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid credentials!"})
		return
	}

	user_token, err := th.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		th.logger.Printf("Error: Authentication CreateNewToken %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error!"})
		return
	}

	utils.WriteJson(w, http.StatusAccepted, utils.Envelope{"auth_token": user_token})
}
