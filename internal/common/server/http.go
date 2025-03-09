package server

import (
	"github.com/Hana-bii/gorder-v2/common/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		// TODO: Warning log
		panic("empty http address")

	}
	RunHTTPServerOnAddr(addr, wrapper)
}

func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	setMiddlewares(apiRouter)
	wrapper(apiRouter)
	apiRouter.Group("/api")
	//apiRouter.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, "pong!!!")
	//})
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}

func setMiddlewares(r *gin.Engine) {
	r.Use(middleware.StructuredLog(logrus.NewEntry(logrus.StandardLogger())))
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLog(logrus.NewEntry(logrus.StandardLogger())))
	r.Use(otelgin.Middleware("default_server"))
}
