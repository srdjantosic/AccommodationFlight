package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"user-service/model"
	"user-service/service"
)

type KeyProduct struct{}
type UserHandler struct {
	Logger  *log.Logger
	Service *service.UserService
}

func NewUserHandler(l *log.Logger, s *service.UserService) *UserHandler {
	return &UserHandler{l, s}
}

func (u *UserHandler) DatabaseName(ctx context.Context) {
	dbs, err := u.Service.Repo.Cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(dbs)
}

func (u *UserHandler) MiddlewareUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		user := &model.User{}
		err := user.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			u.Logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, user)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (u *UserHandler) PostUser(rw http.ResponseWriter, h *http.Request) {
	user := h.Context().Value(KeyProduct{}).(*model.User)
	//newUser := model.User{ID: primitive.NewObjectID(), FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Password: user.Password}
	user.ID = primitive.NewObjectID()
	createdUser, err := u.Service.Insert(user)
	if createdUser == nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
	rw.WriteHeader(http.StatusCreated)
}

func (u *UserHandler) GetUserByEmailAndPassword(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	email := vars["email"]
	password := vars["password"]

	user, err := u.Service.GetUserByEmailAndPassword(email, password)

	if err != nil {
		fmt.Println("Error while logging in.")
		rw.WriteHeader(http.StatusBadRequest)
	}

	if user != nil {
		rw.WriteHeader(http.StatusOK)
	}

}

func (u *UserHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		u.Logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
