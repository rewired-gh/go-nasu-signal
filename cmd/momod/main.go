package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rewired-gh/go-signal-server/internal/app"
	"github.com/rewired-gh/go-signal-server/internal/util"
)

func main() {
	config, err := util.LoadConfig("./")
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	corConfig := cors.DefaultConfig()
	corConfig.AllowWildcard = true
	corConfig.AllowOrigins = []string{"https://*.rewired.moe"}
	server.Use(cors.Default())

	app.HandleServer(server)

	println("Listening on " + config.Listen)
	if config.CertPath != "" && config.KeyPath != "" {
		server.RunTLS(config.Listen, config.CertPath, config.KeyPath)
	} else {
		server.Run(config.Listen)
	}
}
