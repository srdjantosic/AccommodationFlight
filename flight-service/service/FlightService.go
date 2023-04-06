package service

import (
	"flight-service/model"
	"flight-service/repository"
	"log"
)

type FlightService struct {
	Logger *log.Logger
	Repo   *repository.FlightRepository
}

func NewFlightService(l *log.Logger, s *repository.FlightRepository) *FlightService {
	return &FlightService{l, s}
}
func (fs *FlightService) Insert(newFlight *model.Flight) (*model.Flight, error) {

	return fs.Repo.Insert(newFlight) //newUser
}

func (fs *FlightService) Delete(id string) error {

	return fs.Repo.Delete(id)
}

func (fs *FlightService) GetFlightById(id string) (*model.Flight, error) {
	return fs.Repo.GetFlightById(id)
}

func (fs *FlightService) GetFlights(departure string, departurePlace string, arrivalPlace string, noOfSeats int) (model.Flights, error) {
	return fs.Repo.GetAll(departure, departurePlace, arrivalPlace, noOfSeats)
}

func (fs *FlightService) CheckNumberOfFreeSeats(id string, numOfTickets uint64) error {
	flight, err := fs.Repo.GetFlightById(id)
	if err != nil {
		return err
	}
	if numOfTickets > flight.NumberOfFreeSeats {
		return err
	}
	return nil
}

func (fs *FlightService) Update(id string, numberOfTickets uint64) error {

	flight, err := fs.GetFlightById(id)

	if err != nil {
		fs.Logger.Println(err)
		return err
	}

	newNumOfFreeSeats := flight.NumberOfFreeSeats - numberOfTickets

	return fs.Repo.Update(id, newNumOfFreeSeats)
}
