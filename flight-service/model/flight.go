package model

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Flight struct {
	ID                primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Departure         string             `bson:"departure" json:"departure"`
	DeparturePlace    string             `bson:"departurePlace" json:"departurePlace"`
	ArrivalPlace      string             `bson:"arrivalPlace" json:"arrivalPlace"`
	Price             uint64             `bson:"price" json:"price"`
	NumberOfFreeSeats uint64             `bson:"numberOfFreeSeats" json:"numberOfFreeSeats"`
}

type Flights []*Flight

func (f *Flights) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}

func (u *Flight) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}
func (u *Flight) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}
