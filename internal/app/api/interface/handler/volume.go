package handler

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/gin-gonic/gin"
)

type VolumeHandler interface {
	Create(*gin.Context)
}

type volumeHandler struct {
	volumeUC usecase.VolumeUsecase
}

func NewVolumeHandler(volumeUC usecase.VolumeUsecase) VolumeHandler {
	return &volumeHandler{
		volumeUC: volumeUC,
	}
}

func (h *volumeHandler) Create(c *gin.Context) {}
