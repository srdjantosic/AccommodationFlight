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
	"time"
	"user-service/model"
)

type UserRepository struct {
	Cli    *mongo.Client
	Logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*UserRepository, error) {
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
	return &UserRepository{
		Cli:    client,
		Logger: logger,
	}, nil
}
func (u *UserRepository) Disconnect(ctx context.Context) error {
	err := u.Cli.Disconnect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *UserRepository) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.Cli.Ping(ctx, readpref.Primary())
	if err != nil {
		u.Logger.Println(err)
	}

	dbs, err := u.Cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		u.Logger.Println(err)
	}
	fmt.Println(dbs)
}

func (ur *UserRepository) getCollection() *mongo.Collection {
	bookingDatabase := ur.Cli.Database("booking")
	usersCollection := bookingDatabase.Collection("users")
	return usersCollection
}

func (ur *UserRepository) GetByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := ur.getCollection()

	var user model.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		ur.Logger.Println(err)
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByEmailAndPassword(email string, password string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := ur.getCollection()

	var user model.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user)
	if err != nil {
		ur.Logger.Println(err)
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) Insert(user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usersCollection := ur.getCollection()

	result, err := usersCollection.InsertOne(ctx, &user)
	if err != nil {
		ur.Logger.Println(err)
		return nil, err
	}
	ur.Logger.Printf("Documents ID: %v\n", result.InsertedID)
	return user, nil
}
