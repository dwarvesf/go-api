package main

import (
	"github.com/dwarvesf/go-api/docs"
	"github.com/dwarvesf/go-api/pkg/handler"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func setupRouter(a App, sClient *sentry.Client) *gin.Engine {
	docs.SwaggerInfo.Title = "Swagger API"
	docs.SwaggerInfo.Description = "This is a swagger for API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = a.cfg.BaseURL

	if a.cfg.Env == "local" {
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	r := gin.New()

	r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
			AllowHeaders: []string{"Origin", "Host",
				"Content-Type", "Content-Length",
				"Accept-Encoding", "Accept-Language", "Accept",
				"X-CSRF-Token", "Authorization", "X-Requested-With", "X-Access-Token"},
			ExposeHeaders:    []string{"MeAllowMethodsntent-Length"},
			AllowCredentials: true,
		},
	))

	if a.cfg.SentryDSN != "" {
		// Once it's done, you can attach the handler as one of your middleware
		r.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	// handlers
	publicHandler(r, a)
	authenticatedHandler(r, a)

	return r
}

func publicHandler(r *gin.Engine, a App) {
	h := handler.New(a.cfg)

	r.GET("/healthz", h.Healthz)

	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func authenticatedHandler(r *gin.Engine, a App) {

}
