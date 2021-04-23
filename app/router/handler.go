package router

import (
	"github.com/DipandaAser/linker-discord/app"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Start() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", root)
	router.POST("/linker", transfertMessage)

	log.Println("HTTP Server Started on port ", app.Config.HTTPPort)
	err := router.Run(":" + app.Config.HTTPPort)
	log.Fatal(err)
}

func root(c *gin.Context) {
	c.String(http.StatusOK, app.Config.ProjectName)
}

func transfertMessage(c *gin.Context) {
	c.String(http.StatusOK, app.Config.ServiceName)
}

func GetServiceUrl() string {
	return app.Config.WebUrl + "/linker"
}
