package handler

import (
	"net/http"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/builder"
	hdlerr "github.com/atsumarukun/holos-storage-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/pkg/parameter"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
)

type VolumeHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	GetOne(*gin.Context)
	GetAll(*gin.Context)
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
	const errMessage = "failed to create volume"

	var req schema.CreateVolumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeBadRequest, errMessage))
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, errMessage))
		return
	}

	ctx := c.Request.Context()

	volume, err := h.volumeUC.Create(ctx, accountID, req.Name, req.IsPublic)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusCreated, builder.ToVolumeResponse(volume))
}

func (h *volumeHandler) Update(c *gin.Context) {
	const errMessage = "failed to update volume"

	var req schema.UpdateVolumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeBadRequest, errMessage))
		return
	}

	name := c.Param("name")

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, errMessage))
		return
	}

	ctx := c.Request.Context()

	volume, err := h.volumeUC.Update(ctx, accountID, name, req.Name, req.IsPublic)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToVolumeResponse(volume))
}

func (h *volumeHandler) Delete(c *gin.Context) {
	name := c.Param("name")

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to delete volume"))
		return
	}

	ctx := c.Request.Context()

	if err := h.volumeUC.Delete(ctx, accountID, name); err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *volumeHandler) GetOne(c *gin.Context) {
	name := c.Param("name")

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to get volume"))
		return
	}

	ctx := c.Request.Context()

	volume, err := h.volumeUC.GetOne(ctx, accountID, name)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToVolumeResponse(volume))
}

func (h *volumeHandler) GetAll(c *gin.Context) {
	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to get volumes"))
		return
	}

	ctx := c.Request.Context()

	volumes, err := h.volumeUC.GetAll(ctx, accountID)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string][]*schema.VolumeResponse{"volumes": builder.ToVolumeResponses(volumes)})
}
