package app

import (
	"github.com/DipandaAser/linker"
	"os"
)

var Config *linker.ProjectSettings
var DiscordBotToken string

func Init() {

	DiscordBotToken = os.Getenv("BOT_TOKEN")

	Config = &linker.ProjectSettings{}
	Config.ServiceName = "discord"
	Config.ProjectName = "Linker Discord"
	Config.AuthKey = os.Getenv("APIKEY")
	Config.DBName = os.Getenv("DB_NAME")
	Config.MongodbURI = os.Getenv("MONGO_URI")
	Config.HTTPPort = os.Getenv("HTTP_PORT")
	Config.WebUrl = os.Getenv("WEB_URL")
}
