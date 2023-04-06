package handler

import (
	"context"
	"flight-service/model"
	"flight-service/service"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
)

type KeyProduct struct{}
type FlightHandler struct {
	Logger  *log.Logger
	Service *service.FlightService
}

func NewFlightHandler(l *log.Logger, s *service.FlightService) *FlightHandler {
	return &FlightHandler{l, s}
}

func (fh *FlightHandler) DatabaseName(ctx context.Context) {
	dbs, err := fh.Service.Repo.Cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(dbs)
}

func (fh *FlightHandler) MiddlewareUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		flight := &model.Flight{}
		err := flight.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			fh.Logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, flight)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (fh *FlightHandler) PostFlight(rw http.ResponseWriter, h *http.Request) {
	flight := h.Context().Value(KeyProduct{}).(*model.Flight)
	//newUser := model.User{ID: primitive.NewObjectID(), FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Password: user.Password}
	flight.ID = primitive.NewObjectID()
	createdFlight, err := fh.Service.Insert(flight)
	if createdFlight == nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
	rw.WriteHeader(http.StatusCreated)
}

func (fh *FlightHandler) GetFlights(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	date := vars["departure"]
	departurePlace := vars["departurePlace"]
	arrivalPlace := vars["arrivalPlace"]
	noOfSeats := vars["noOfSeats"]

	n, err := strconv.Atoi(noOfSeats)
	if err != nil {
		fmt.Println("Error during conversion.")
		return
	}
	fmt.Println("Broj iz URLA: ", n)

	flights, err := fh.Service.GetFlights(date, departurePlace, arrivalPlace, n)

	if err != nil {
		fh.Logger.Println("Database exception ", err)
	}

	if flights == nil {
		return
	}
	err = flights.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		fh.Logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (fh *FlightHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		fh.Logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

func (fh *FlightHandler) DeleteFlight(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	fh.Service.Delete(id)
	rw.WriteHeader(http.StatusNoContent)
}

func (fh *FlightHandler) GetFlightById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	flight, err := fh.Service.GetFlightById(id)
	if err != nil {
		fh.Logger.Print("Database exception: ", err)
	}

	if flight == nil {
		http.Error(rw, "Flight with given id not found", http.StatusNotFound)
		fh.Logger.Printf("Flight with id: '%s' not found", id)
		return
	}

	err = flight.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		fh.Logger.Fatal("Unable to convert to json :", err)
		return
	}

}

func (fh *FlightHandler) CheckNumberOfFreeSeats(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	flightId := vars["flightId"]
	numOfTickets := vars["numberOfTickets"]

	numberOfTickets, err := strconv.ParseUint(numOfTickets, 10, 64)

	if err != nil {
		fh.Logger.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
	}
	err = fh.Service.CheckNumberOfFreeSeats(flightId, numberOfTickets)
	if err != nil {
		fh.Logger.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
	}
	rw.WriteHeader(http.StatusOK)
}

func (fh *FlightHandler) Update(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	flightId := vars["flightId"]
	numOfTickets := vars["numberOfTickets"]
	numberOfTickets, err := strconv.ParseUint(numOfTickets, 10, 64)

	if err != nil {
		fh.Logger.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
	}

	err = fh.Service.Update(flightId, numberOfTickets)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
	}
	rw.WriteHeader(http.StatusOK)
}
