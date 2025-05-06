package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/builder"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/pkg/parameter"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *entryHandler) Create(c *gin.Context) {
	var req schema.CreateEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err)
		errors.Handle(c, status.Error(code.BadRequest, "failed to parse multipart/form-data"))
		return
	}

	isPublic, err := strconv.ParseBool(req.IsPublic)
	if err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "failed to parse is_public to bool"))
		return
	}

	volumeID, err := uuid.Parse(req.VolumeID)
	if err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "failed to parse volume_id to uuid"))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "failed to get file"))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "failed to open file"))
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			errors.Handle(c, err)
		}
	}()

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	entry, err := h.entryUC.Create(ctx, accountID, volumeID, req.Key, uint64(fileHeader.Size), isPublic, file)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusCreated, builder.ToEntryResponse(entry))
}
