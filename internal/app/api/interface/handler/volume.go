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
	Update(*gin.Context)
	Delete(*gin.Context)
	GetOne(*gin.Context)
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
		errors.Handle(c, status.Error(code.BadRequest, "failed to parse json"))
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

func (h *volumeHandler) Update(c *gin.Context) {
	var req schema.UpdateVolumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "failed to parse json"))
		return
	}

	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "invalid id"))
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	volume, err := h.volumeUC.Update(ctx, accountID, id, req.Name, req.IsPublic)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToVolumeResponse(volume))
}

func (h *volumeHandler) Delete(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "invalid id"))
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	if err := h.volumeUC.Delete(ctx, accountID, id); err != nil {
		errors.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *volumeHandler) GetOne(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "invalid id"))
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	volume, err := h.volumeUC.GetOne(ctx, accountID, id)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToVolumeResponse(volume))
}
