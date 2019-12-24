package logformatter_test

import (
	"github.com/gin-gonic/gin"
	f "github.com/twreporter/logformatter"
)

func ExampleNewGinLogFormatter() {
	// disable default logger(stdout/stderr)
	r := gin.New()

	r.Use(gin.Recovery())

	// customize log severity here
	// default to Info
	formatter := f.NewGinLogFormatter(f.GinLogSeverity(f.Debug))

	// config gin with the customize logger
	r.Use(gin.LoggerWithFormatter(formatter))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "{message: \"pong\"}")
	})

	r.Run(":8080")
}
