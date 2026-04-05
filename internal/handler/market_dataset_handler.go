package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tradingsystem/internal/dto"
	"tradingsystem/internal/service"
	"tradingsystem/pkg/response"
)

type MarketDatasetHandler struct {
	service *service.MarketDatasetService
}

func NewMarketDatasetHandler(service *service.MarketDatasetService) *MarketDatasetHandler {
	return &MarketDatasetHandler{service: service}
}

func (h *MarketDatasetHandler) Import(c *gin.Context) {
	symbol := c.PostForm("symbol")
	timeframe := c.PostForm("timeframe")
	if symbol == "" || timeframe == "" {
		response.Error(c, http.StatusBadRequest, "symbol and timeframe are required")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file is required")
		return
	}

	data, err := h.service.Import(symbol, timeframe, fileHeader)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *MarketDatasetHandler) List(c *gin.Context) {
	var req dto.ListMarketDatasetsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.service.List(req.Symbol, req.Timeframe)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"items": items})
}

func (h *MarketDatasetHandler) Get(c *gin.Context) {
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

func (h *MarketDatasetHandler) Delete(c *gin.Context) {
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
