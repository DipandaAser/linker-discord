package main

import (
	"context"
	"github.com/DipandaAser/linker"
	"github.com/DipandaAser/linker-discord/app"
	"github.com/DipandaAser/linker-discord/app/router"
	"github.com/DipandaAser/linker-discord/bot"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

var db *mongo.Database
var ctx context.Context = context.TODO()

func main() {

	_ = godotenv.Load()
	app.Init()
	linker.MongoCtx = &ctx

	// ─── MONGO ──────────────────────────────────────────────────────────────────────
	err := MongoConnect()
	if err != nil {
		log.Fatal("Can't setup mongodb")
	}

	// ─── WE REFRESH THE MONGO CONNECTION EACH 10MINS ──────────────────────────────────────
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			go MongoReconnectCheck()
		}
	}()

	go router.Start()

	// we put te service online to allow other service to link
	_, err = linker.SetService(app.Config.ServiceName, router.GetServiceUrl(), app.Config.AuthKey, linker.StatusOnline)
	if err != nil {
		log.Fatal("can't set service")
	}

	// We put the service offline if the program stop
	defer linker.SetService(app.Config.ServiceName, router.GetServiceUrl(), app.Config.AuthKey, linker.StatusOffline)

	err = bot.Init()
	if err != nil {
		log.Fatal(err)
	}
}

// MongoConnect connects to mongoDB
func MongoConnect() error {

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// We make sure we have been connected
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	db = client.Database(os.Getenv("DB_NAME"))
	linker.DB = db

	return nil
}

// MongoReconnectCheck reconnects to MongoDB
func MongoReconnectCheck() {

	// We make sure we are still connected
	err := db.Client().Ping(ctx, readpref.Primary())
	if err != nil {
		// We reconnect
		_ = MongoConnect()
	}
}
