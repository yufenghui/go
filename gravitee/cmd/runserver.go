package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/yufenghui/go/gravitee/services"
	"strconv"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

func RunServer(configFile string) error {

	cfg, db, err := initConfig(configFile)
	if err != nil {
		return err
	}
	defer db.Close()

	// start the services
	if err := services.Init(cfg, db); err != nil {
		return err
	}
	defer services.Close()

	// create gin app
	g := gin.Default()
	loadRouter(g)

	graceful.Run(":"+strconv.Itoa(cfg.ServerPort), 5*time.Second, g)

	return nil
}

func loadRouter(g *gin.Engine) {

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	service := g.Group("/v1")
	{
		service.GET("/", services.HealthService.HealthCheck)
	}

}
