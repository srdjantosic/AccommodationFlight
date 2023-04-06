package main

import (
	"context"
	"flight-service/handler"
	"flight-service/repository"
	"flight-service/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	flightLogger := log.New(os.Stdout, "[flight-store] ", log.LstdFlags)

	flightStore, err := repository.New(timeoutContext, flightLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer flightStore.Disconnect(timeoutContext)

	flightStore.Ping()
	flightService := service.NewFlightService(logger, flightStore)
	flightsHandler := handler.NewFlightHandler(logger, flightService)

	flightsHandler.DatabaseName(timeoutContext)

	router := mux.NewRouter()
	router.Use(flightsHandler.MiddlewareContentTypeSet)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/{flightId}/{numberOfTickets}", flightsHandler.CheckNumberOfFreeSeats)

	updateRouter := router.Methods(http.MethodGet).Subrouter()
	updateRouter.HandleFunc("/update/{flightId}/{numberOfTickets}", flightsHandler.Update)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", flightsHandler.PostFlight)
	postRouter.Use(flightsHandler.MiddlewareUserDeserialization)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/delete/{id}", flightsHandler.DeleteFlight)

	getFlightByIdRouter := router.Methods(http.MethodGet).Subrouter()
	getFlightByIdRouter.HandleFunc("/{id}", flightsHandler.GetFlightById)

	getFlights := router.Methods(http.MethodGet).Subrouter()
	getFlights.HandleFunc("/{departure}/{departurePlace}/{arrivalPlace}/{noOfSeats}", flightsHandler.GetFlights)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")

}
