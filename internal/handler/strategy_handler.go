package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tradingsystem/internal/dto"
	"tradingsystem/internal/service"
	"tradingsystem/pkg/response"
)

type StrategyHandler struct {
	service *service.StrategyService
}

func NewStrategyHandler(service *service.StrategyService) *StrategyHandler {
	return &StrategyHandler{service: service}
}

func (h *StrategyHandler) Create(c *gin.Context) {
	var req dto.CreateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *StrategyHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"items": items})
}

func (h *StrategyHandler) Get(c *gin.Context) {
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

func (h *StrategyHandler) Update(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.Update(id, req)
	if err != nil {
		handleRepositoryError(c, err)
		return
	}

	response.Success(c, data)
}

func (h *StrategyHandler) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		handleRepositoryError(c, err)
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func parseUintParam(c *gin.Context, name string) (uint, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}

	return uint(id), true
}

func handleRepositoryError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Error(c, http.StatusNotFound, "record not found")
		return
	}

	response.Error(c, http.StatusInternalServerError, err.Error())
}
