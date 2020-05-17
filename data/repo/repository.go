package repo

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"net/http"
	. "rafiki/data/model"
	. "rafiki/engine/sms"
)

func SendMessageHandler(c *gin.Context) {

	database := c.MustGet("database").(*gorm.DB)

	var payload Message
	if err := c.ShouldBindJSON(&payload); err != nil {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				hub.CaptureException(err)
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := Message{
		MessageId:       uuid.New(),
		UserId:          payload.UserId,
		RecipientNumber: payload.RecipientNumber,
		MessageBody:     payload.MessageBody,
	}

	database.Create(&message)

	go sendMessage(message, database)

	c.JSON(http.StatusOK, message)
}

func sendMessage(message Message, database *gorm.DB) {

	response, status := SendATMessage(
		message.RecipientNumber,
		message.MessageBody,
	)

	if status == true {
		var messageId = message.MessageId.String()
		var message Message
		if err := database.Where("message_id = ?", messageId).First(&message).Error; err != nil {
			return
		}

		database.Model(&message).Update("response", response)
	} else {
		var messageId = message.MessageId.String()
		var message Message
		if err := database.Where("message_id = ?", messageId).First(&message).Error; err != nil {
			return
		}

		database.Model(&message).Update("message_sent", false)
	}
}

func FetchMessageHandler(c *gin.Context) {

	database := c.MustGet("database").(*gorm.DB)

	var message Message
	if err := database.Where("message_id = ?", c.Param("message_id")).First(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, message)
}

func FetchAllMessagesHandler(c *gin.Context) {

	database := c.MustGet("database").(*gorm.DB)

	var message []Message
	database.Find(&message)

	c.JSON(http.StatusOK, message)
}
