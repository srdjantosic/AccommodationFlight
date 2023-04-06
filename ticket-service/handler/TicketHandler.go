package handler

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"ticket-service/model"
	"ticket-service/service"
)

type KeyProduct struct{}

type TicketHandler struct {
	Logger  *log.Logger
	Service *service.TicketService
}

func NewTicketHandler(l *log.Logger, s *service.TicketService) *TicketHandler {
	return &TicketHandler{l, s}
}

func (th *TicketHandler) DatabaseName(ctx context.Context) {
	dbs, err := th.Service.Repo.Cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(dbs)
}

func (th *TicketHandler) MiddlewareUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		ticket := &model.Ticket{}
		err := ticket.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			th.Logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, ticket)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (th *TicketHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		th.Logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

func (th *TicketHandler) PostTicket(rw http.ResponseWriter, h *http.Request) {

	ticket := h.Context().Value(KeyProduct{}).(*model.Ticket)
	ticket.ID = primitive.NewObjectID()

	createdTicket, err := th.Service.Insert(ticket)
	if createdTicket == nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {

		rw.WriteHeader(http.StatusBadRequest)
	}
	rw.WriteHeader(http.StatusCreated)
}

func (th *TicketHandler) GetTicketsByUserId(rw http.ResponseWriter, h *http.Request) {
	userId := h.URL.Query().Get("userId")

	tickets, err := th.Service.GetByUserId(userId)
	if err != nil {
		th.Logger.Println("Database exception: ", err)
	}
	if tickets == nil {
		return
	}
	err = tickets.ToJson(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		th.Logger.Fatal("Unable to convert to json: ", err)
		return
	}
}
