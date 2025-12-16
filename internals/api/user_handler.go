package api

import (
	"backend-go/internals/store"
	"backend-go/internals/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) userRegisterValidation(user *registerUserRequest) error {
	if user.Username == "" {
		return errors.New("username is required")
	}

	if len(user.Username) > 50 {
		return errors.New("username cannot be more than 50 characters")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("email is not valid")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userRequest := &registerUserRequest{}
	err := json.NewDecoder(r.Body).Decode(userRequest)
	if err != nil {
		uh.logger.Printf("Error: CreateUser Decode %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	err = uh.userRegisterValidation(userRequest)
	if err != nil {
		uh.logger.Printf("Error: Not valid request%v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		uh.logger.Printf("Error: Bcrypting %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "User is already exist"})
		return
	}

	user := &store.User{
		Email:    userRequest.Email,
		Username: userRequest.Username,
		Password: string(hashedPassword),
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		if err == store.ErrEmailExisted {
			uh.logger.Printf("Error: %v\n", err)
			utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Email Address is already exist"})
			return
		}
		if err == store.ErrUsernameExisted {
			uh.logger.Printf("Error: %v\n", err)
			utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Username is already exist"})
			return
		}
		uh.logger.Printf("Error: CreateUser Store%v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"data": user})
}
