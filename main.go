package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"rafiki/data/repo"
	"rafiki/health"
	. "rafiki/settings"
	"strings"
	"syscall"
	"time"
)

func setupRouter() *gin.Engine {

	// Force log's color
	gin.ForceConsoleColor()

	// Creates a gin router with default middleware:
	// - logger and recovery (crash-free) middleware
	router := gin.Default()

	//set the gin mode
	gin.SetMode(gin.DebugMode)

	router.Use(cors.Default())

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	return router
}

func main() {

	// logs will write with UNIX time
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	router := setupRouter()

	// -----------------------------------------------------------------------------------------------------------------
	// ROOT API ENDPOINT
	// -----------------------------------------------------------------------------------------------------------------

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "Rafiki service"})
	})

	// -----------------------------------------------------------------------------------------------------------------
	// SERVICE HEALTH API
	// -----------------------------------------------------------------------------------------------------------------

	router.GET("/api/health", health.CheckHealthHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// MESSAGE API's
	// -----------------------------------------------------------------------------------------------------------------

	router.POST("/api/v1/message/send", repo.SendMessageHandler)

	router.GET("/api/v1/message/send/:message_id", repo.FetchMessageHandler)

	router.GET("/api/v1/message/send", repo.FetchAllMessagesHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// START THE GIN SERVER
	// -----------------------------------------------------------------------------------------------------------------

	srv := &http.Server{
		Addr:    getPort(),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msgf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Print("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msgf("Server forced to shutdown : %v", err)
	}

	logger.Print("Server exiting")

}

func getPort() string {
	port := os.Getenv("PORT")
	if GetEnv() == "DEMO" {
		if port == "" {
			port = "8080"
			log.Print("no DEMO PORT environment variable detected. Setting port to ", port)
		}
	} else if GetEnv() == "STAGING" {
		if port == "" {
			port = "8080"
			log.Print("no STAGING PORT environment variable detected. Setting port to ", port)
		}
	} else if GetEnv() == "PRODUCTION" {
		if port == "" {
			port = "8080"
			log.Print("no PRODUCTION PORT environment variable detected. Setting port to ", port)
		}
	}
	return ":" + port
}
