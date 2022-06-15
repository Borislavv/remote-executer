package remoter

import (
	"context"
	"time"

	"github.com/serge64/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	agg "github.com/Borislavv/remote-executer/internal/domain/agg/msg"
	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/domain/usecase"
	"github.com/Borislavv/remote-executer/internal/util"
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
		return util.ErrWithTrace(err)
	}

	// init. context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect to mongodb
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return util.ErrWithTrace(err)
	}
	defer mongoClient.Disconnect(ctx)

	// check the db is available
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return util.ErrWithTrace(err)
	}

	// repos.
	msgRepo := mongoRepo.NewMsgRepo(
		mongoClient.Database(config.MongoDB).Collection(MongoMsgsCollection),
	)

	// app deps.
	gateway := usecase.NewTelegram(config.TelegramEndpoint, config.TelegramToken)

	// usecases
	polling := usecase.NewPolling(ctx, gateway, msgRepo)
	messages := usecase.NewMessages(ctx, msgRepo)
	commands := usecase.NewCommands(ctx, msgRepo)
	responses := usecase.NewResponses(ctx, gateway)

	// chans
	messagesCh := make(chan []agg.Msg)
	responseCh := make(chan dto.TelegramResponseInterface)
	errCh := make(chan error)

	go polling.Do(messagesCh, errCh)
	go messages.Consuming(messagesCh, errCh)
	go commands.Executing(responseCh, errCh)
	go responses.Sending(responseCh, errCh)

	for {
		select {
		case err := <-errCh:
			cancel()

			// awaiting while all gorutines will finished
			time.Sleep(time.Second)

			return err
		default:
			// don't block waiting for an error
		}
	}
}
