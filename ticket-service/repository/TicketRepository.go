package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"ticket-service/model"
	"time"
)

type TicketRepository struct {
	Cli    *mongo.Client
	Logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*TicketRepository, error) {

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
	return &TicketRepository{
		Cli:    client,
		Logger: logger,
	}, nil
}
func (tr *TicketRepository) Disconnect(ctx context.Context) error {
	err := tr.Cli.Disconnect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (tr *TicketRepository) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := tr.Cli.Ping(ctx, readpref.Primary())
	if err != nil {
		tr.Logger.Println(err)
	}

	dbs, err := tr.Cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		tr.Logger.Println(err)
	}
	fmt.Println(dbs)
}

func (tr *TicketRepository) getCollection() *mongo.Collection {
	return tr.Cli.Database("booking").Collection("tickets")
}

func (tr *TicketRepository) Insert(ticket *model.Ticket) (*model.Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usersCollection := tr.getCollection()

	result, err := usersCollection.InsertOne(ctx, &ticket)
	if err != nil {
		tr.Logger.Println(err)
		return nil, err
	}
	tr.Logger.Printf("Documents ID: %v\n", result.InsertedID)
	return ticket, nil
}

func (tr *TicketRepository) GetByUserId(id string) (model.Tickets, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticketsCollection := tr.getCollection()

	var tickets model.Tickets
	ticketsCursor, err := ticketsCollection.Find(ctx, bson.M{"userId": id})
	if err != nil {
		tr.Logger.Println(err)
		return nil, err
	}
	if err = ticketsCursor.All(ctx, &tickets); err != nil {
		tr.Logger.Println(err)
		return nil, err
	}
	return tickets, nil
}
