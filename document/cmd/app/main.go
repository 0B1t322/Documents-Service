package main

import (
	"fmt"
	"github.com/0B1t322/Documents-Service/document/internal/app"
	"github.com/0B1t322/Documents-Service/document/internal/config"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
		)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	if err := config.FromEnv(); err != nil {
		panic(err)
	}
	cfg := config.GlobalConfig

	app, err := app.NewHTTPAppFromConfig(cfg)
	if err != nil {
		panic(err)
	}

	eng := gin.New()

	api := eng.Group("/api")
	{
		api.Use(CORSMiddleware())
		documents := api.Group("documents")
		{
			v1 := documents.Group("/v1")
			{
				appHandler, err := app.ToHandler(v1.BasePath())
				if err != nil {
					panic(err)
				}

				v1.Any(
					"/*any", func(c *gin.Context) {
						appHandler.ServeHTTP(c.Writer, c.Request)
					},
				)
			}

			swagger := documents.Group("/swagger")
			{
				swagger.Static("", "./api/open-api")
			}
		}
	}

	if err := eng.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		panic(err)
	}
}
