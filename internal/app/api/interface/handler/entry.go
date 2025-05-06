package handler

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/gin-gonic/gin"
)

type EntryHandler interface {
	Create(*gin.Context)
}

type entryHandler struct {
	entryUC usecase.EntryUsecase
}

func NewEntryHandler(entryUC usecase.EntryUsecase) EntryHandler {
	return &entryHandler{
		entryUC: entryUC,
	}
}

func (h *entryHandler) Create(c *gin.Context) {}
