package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	stock "github.com/alireza-aliabadi/golang-stocking-update/internal/models"
	rabbitmq "github.com/alireza-aliabadi/golang-stocking-update/internal/rabbitmq"
)

type StockHandler struct {
	Rabbitmq *rabbitmq.RabbitClient
}

func NewStockHandler(rc *rabbitmq.RabbitClient) *StockHandler {
	return &StockHandler{
		Rabbitmq: rc,
	}
}

func (h *StockHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/stock", h.UpdateStock)
}

func (h* StockHandler) UpdateStock(c echo.Context) error {
	var payload stock.StockUpdatePayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := h.Rabbitmq.Publish(payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to queue update"})
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "Stock update queued successfully",
		"stock-unit": payload.StockUnit,
	})
}