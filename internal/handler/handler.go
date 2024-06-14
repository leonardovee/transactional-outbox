package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Command(c *gin.Context) {
	ctx := c.Request.Context()

	var request struct {
		AggregateID *string `json:"AggregateID"`
		Command     string  `json:"Command"`
		Total       *int    `json:"Total"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := h.service.Execute(ctx, &service.Command{
		Command:     request.Command,
		AggregateID: request.AggregateID,
		Total:       request.Total,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, o)
}
