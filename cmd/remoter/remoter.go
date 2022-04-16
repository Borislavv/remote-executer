package remoter

import (
	"context"

	"github.com/serge64/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoUsersCollection           = "users"
	MongoCommandsHistoryCollection = "commandsHistory"
)

type Config struct {
	// store
	MongoURI string `env:"MONGO_URI,default=mongodb://localhost:27017/"`
	MongoDB  string `env:"MONGO_DATABASE,default=remoter"`

	// service props.
	WorkerTimeout     int `env:"WORKER_TIMEOUT,default=1"`
	WriteMongoTimeout int `env:"WRITE_MONGO_TIMEOUT,default=1"`
	ReadMongoTimeout  int `env:"READ_MONGO_TIMEOUT,default=1"`
}

func Run() error {
	// init. app config
	config := Config{}
	if err := env.Unmarshal(&config); err != nil {
		return err
	}

	// init. context
	ctx := context.Background()
	cancelCtx, cancelClosure := context.WithCancel(ctx)
	defer cancelClosure()

	// connect to mongodb
	mongoClient, err := mongo.Connect(cancelCtx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}
	defer mongoClient.Disconnect(cancelCtx)

	// check the db is available
	if err := mongoClient.Ping(cancelCtx, readpref.Primary()); err != nil {
		return err
	}

	return nil
}
