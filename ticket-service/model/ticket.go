package model

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Ticket struct {
	ID              primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"userId, omitempty" json:"userId"`
	FlightID        primitive.ObjectID `bson:"flightId, omitempty" json:"flightId"`
	NumberOfTickets uint8              `bson:"numberOfTickets" json:"numberOfTickets"`
}

type Tickets []*Ticket

func (t *Tickets) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (u *Ticket) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}
func (u *Ticket) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}
