package logformatter_test

import (
	"os"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	f "github.com/twreporter/logformatter"
	"google.golang.org/genproto/googleapis/logging/type"
)

func ExampleNewGinLogFormatter() {
	// disable default logger(stdout/stderr)
	r := gin.New()

	r.Use(gin.Recovery())

	// customize log severity here
	// default to Info
	formatter := f.NewGinLogFormatter(f.GinLogSeverity(ltype.LogSeverity_DEBUG))

	// config gin with the customize logger
	r.Use(gin.LoggerWithFormatter(formatter))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "{message: \"pong\"}")
	})

	r.Run(":8080")
}

func ExampleNewStackdriverFormatter() {
	log.SetOutput(os.Stdout)

	// Config stackdriver formatter for example service
	log.SetFormatter(f.NewStackdriverFormatter("example", "test"))

	// Log message with info severity
	log.Info("message")

	// Format the error objecgt in runtime.Stack shape
	log.Error(f.FormatStack(errors.Errorf("error")))
}
