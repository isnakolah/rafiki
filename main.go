package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"recoin-notification-service/data/repo"
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
	if settings.GetEnv() == "DEMO" {
		if port == "" {
			port = "4003"
			log.Print("no DEMO PORT environment variable detected. Setting port to ", port)
		}
	} else if settings.GetEnv() == "STAGING" {
		if port == "" {
			port = "4003"
			log.Print("no STAGING PORT environment variable detected. Setting port to ", port)
		}
	} else if settings.GetEnv() == "PRODUCTION" {
		if port == "" {
			port = "4003"
			log.Print("no PRODUCTION PORT environment variable detected. Setting port to ", port)
		}
	}
	return ":" + port
}
