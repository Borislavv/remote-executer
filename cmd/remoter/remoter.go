package remoter

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/serge64/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	"github.com/Borislavv/remote-executer/internal/domain/agg"
	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/domain/errs"
	"github.com/Borislavv/remote-executer/internal/domain/usecase"
)

const (
	MongoUsersCollection = "users"
	MongoMsgsCollection  = "messages"
	MongoChatsCollection = "chats"
)

type Config struct {
	// store
	MongoURI string `env:"MONGO_URI,default=mongodb://localhost:27017/"`
	MongoDB  string `env:"MONGO_DATABASE,default=remoter"`

	TelegramUsername string `env:"TELEFRAM_USERNAME,default=BorislavGlazunov"`
	TelegramEndpoint string `env:"TELEGRAM_ENDPOINT,default=https://api.telegram.org/"`
	TelegramToken    string `env:"TELEGRAM_TOKEN,default="`

	// timeout's in Milliseconds
	PollingTimeout int `env:"POLLING_TIMEOUT,default=450"`
	MongoDbTimeout int `env:"MONGODB_TIMEOUT,default=450"`
}

func Run() error {
	// init. app config
	cfg := Config{}
	if err := env.Unmarshal(&cfg); err != nil {
		return errs.New(err).Interrupt()
	}

	// init. ctx and wg
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect to mongodb
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return errs.New(err).Interrupt()
	}
	defer mongoClient.Disconnect(ctx)

	// check the db is available
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return errs.New(err).Interrupt()
	}

	// repos.
	msgRepo := mongoRepo.NewMsgRepo(
		mongoClient.Database(cfg.MongoDB).Collection(MongoMsgsCollection),
	)

	// app deps.
	gateway := usecase.NewTelegram(cfg.TelegramEndpoint, cfg.TelegramToken)

	// usecases
	polling := usecase.NewPolling(ctx, gateway, msgRepo, wg, cfg.PollingTimeout)
	messages := usecase.NewMessages(ctx, msgRepo, wg)
	commands := usecase.NewCommands(ctx, msgRepo, wg, cfg.MongoDbTimeout, cfg.TelegramUsername)
	responses := usecase.NewResponses(ctx, gateway, wg)

	// channels
	messagesCh := make(chan []agg.Msg)
	responseCh := make(chan dto.TelegramResponseInterface)
	errCh := make(chan error, 1)
	sysSignalsCh := make(chan os.Signal, 1)
	signal.Notify(sysSignalsCh, os.Interrupt)

	wg.Add(4)
	go polling.Do(messagesCh, errCh)
	go messages.Consuming(messagesCh, errCh)
	go commands.Executing(responseCh, errCh)
	go responses.Sending(responseCh, errCh)

	return runApp(wg, cancel, errCh, sysSignalsCh)
}

func runApp(
	wg *sync.WaitGroup,
	cancel context.CancelFunc,
	errCh <-chan error,
	sysSigsCh <-chan os.Signal,
) error {
	defer func() {
		cancel()
		wg.Wait()
	}()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				e, ok := err.(errs.ErrorWithTrace)
				if ok && e.IsInterrupt() {
					return e
				}
				log.Println(err.Error())
			}
		case <-sysSigsCh:
			log.Println("interception the interrupt signal CTRL+C, stopping the app")
			return nil
		}
	}
}
