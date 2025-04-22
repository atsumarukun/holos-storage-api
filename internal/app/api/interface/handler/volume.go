package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/builder"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/pkg/parameter"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
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

func (h *volumeHandler) Create(c *gin.Context) {
	var req schema.CreateVolumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "bad request"))
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	volume, err := h.volumeUC.Create(ctx, accountID, req.Name, req.IsPublic)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusCreated, builder.ToVolumeResponse(volume))
}
