package repository

import (
	"context"
	"flight-service/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type FlightRepository struct {
	Cli    *mongo.Client
	Logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*FlightRepository, error) {

	dburi := os.Getenv("MONGODB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &FlightRepository{
		Cli:    client,
		Logger: logger,
	}, nil
}
func (fr *FlightRepository) Disconnect(ctx context.Context) error {
	err := fr.Cli.Disconnect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (fr *FlightRepository) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := fr.Cli.Ping(ctx, readpref.Primary())
	if err != nil {
		fr.Logger.Println(err)
	}

	dbs, err := fr.Cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		fr.Logger.Println(err)
	}
	fmt.Println(dbs)
}

func (fr *FlightRepository) getCollection() *mongo.Collection {
	bookingDatabase := fr.Cli.Database("booking")
	usersCollection := bookingDatabase.Collection("flights")
	return usersCollection
}

func (fr *FlightRepository) Insert(flight *model.Flight) (*model.Flight, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usersCollection := fr.getCollection()

	result, err := usersCollection.InsertOne(ctx, &flight)
	if err != nil {
		fr.Logger.Println(err)
		return nil, err
	}
	fr.Logger.Printf("Documents ID: %v\n", result.InsertedID)
	return flight, nil
}

func (fr *FlightRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	flightsCollection := fr.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	result, err := flightsCollection.DeleteOne(ctx, filter)
	if err != nil {
		fr.Logger.Println(err)
		return err
	}
	fr.Logger.Printf("Documents deleted: %v\n", result.DeletedCount)
	return nil
}

func (fr *FlightRepository) GetFlightById(id string) (*model.Flight, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	flightsCollection := fr.getCollection()

	fmt.Println(id)

	var flight model.Flight
	objID, _ := primitive.ObjectIDFromHex(id)
	err := flightsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&flight)
	if err != nil {
		fr.Logger.Println(err, "tu sam")
		return nil, err
	}
	return &flight, nil
}

func (fr *FlightRepository) GetAll(departure string, departurePlace string, arrivalPlace string, noOfSeats int) (model.Flights, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	flightsCollection := fr.getCollection()

	var flights model.Flights
	flightsCursor, err := flightsCollection.Find(ctx, bson.M{
		"departure":         departure,
		"departurePlace":    departurePlace,
		"arrivalPlace":      arrivalPlace,
		"numberOfFreeSeats": bson.M{"$gte": noOfSeats}})

	if err != nil {
		fr.Logger.Println(err)
		return nil, err
	}
	if err = flightsCursor.All(ctx, &flights); err != nil {
		fr.Logger.Println(err)
		return nil, err
	}
	return flights, nil
}

func (fr *FlightRepository) Update(id string, newNumberOfFreeSeats uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	flightsCollection := fr.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"numberOfFreeSeats": newNumberOfFreeSeats,
	}}
	result, err := flightsCollection.UpdateOne(ctx, filter, update)
	fr.Logger.Printf("Documents matched: %v\n", result.MatchedCount)
	fr.Logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		fr.Logger.Println(err)
		return err
	}
	return nil
}
