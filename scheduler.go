package main

import (
    "net/http"
	"log"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type SchedulerHandler struct {
    rabbitMQ *RabbitMQService
}

func NewSchedulerHandler(rabbitMQ *RabbitMQService) *SchedulerHandler {
    return &SchedulerHandler{rabbitMQ: rabbitMQ}
}

func (h *SchedulerHandler) ScheduleEmail(c *gin.Context) {
    var email Email
    if err := c.ShouldBindJSON(&email); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    email.ID = uuid.New().String()

    if err := h.rabbitMQ.ScheduleEmail(email); err != nil {
				log.Printf("Error scheduling email: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Email scheduled successfully", "id": email.ID})
}
