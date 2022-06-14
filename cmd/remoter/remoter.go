package remoter

import (
	"context"

	"github.com/serge64/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	"github.com/Borislavv/remote-executer/internal/domain/usecase"
	telegramGateway "github.com/Borislavv/remote-executer/pkg/gateway/telegram"
)

const (
	MongoUsersCollection           = "users"
	MongoMsgsCollection            = "messages"
	MongoCommandsHistoryCollection = "commandsHistory"
)

type Config struct {
	// store
	MongoURI string `env:"MONGO_URI,default=mongodb://localhost:27017/"`
	MongoDB  string `env:"MONGO_DATABASE,default=remoter"`

	TelegramEndpoint string `env:"TELEGRAM_ENDPOINT,default=https://api.telegram.org/"`
	TelegramToken    string `env:"TELEGRAM_TOKEN,default=5022497048:AAGQcUiyExpJXr3pjjv_cgVody3rv_MvjZ4"`

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect to mongodb
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}
	defer mongoClient.Disconnect(ctx)

	// check the db is available
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	// repos.
	msgRepo := mongoRepo.NewMsgRepo(
		mongoClient.Database(config.MongoDB).Collection(MongoMsgsCollection),
	)

	// app deps.
	gateway := usecase.NewTelegram(config.TelegramEndpoint, config.TelegramToken)
	polling := usecase.NewPolling(ctx, gateway, msgRepo)

	// chans
	messagesCh := make(chan telegramGateway.ResponseGetMessagesInterface)
	errCh := make(chan error)

	go polling.Do(messagesCh, errCh)

	for {
		select {
		case err := <-errCh:
			cancel()
			return err
		default:
			// don't block waiting for an error
		}
	}
}
