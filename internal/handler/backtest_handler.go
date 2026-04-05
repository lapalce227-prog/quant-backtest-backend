package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tradingsystem/internal/dto"
	"tradingsystem/internal/service"
	"tradingsystem/pkg/response"
)

type BacktestHandler struct {
	service *service.BacktestService
}

func NewBacktestHandler(service *service.BacktestService) *BacktestHandler {
	return &BacktestHandler{service: service}
}

func (h *BacktestHandler) Create(c *gin.Context) {
	var req dto.CreateBacktestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.Create(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "record not found")
			return
		}

		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *BacktestHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"items": items})
}

func (h *BacktestHandler) Get(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	data, err := h.service.GetByID(id)
	if err != nil {
		handleRepositoryError(c, err)
		return
	}

	response.Success(c, data)
}

func (h *BacktestHandler) ListTrades(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	items, err := h.service.ListTrades(id)
	if err != nil {
		handleRepositoryError(c, err)
		return
	}

	response.Success(c, gin.H{"items": items})
}
