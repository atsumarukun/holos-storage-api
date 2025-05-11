package handler

import (
	errs "errors"
	"mime/multipart"
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

type EntryHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
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
		errors.Handle(c, status.Error(code.BadRequest, "failed to parse multipart/form-data"))
		return
	}

	volumeID, err := uuid.Parse(req.VolumeID)
	if err != nil {
		errors.Handle(c, status.Error(code.BadRequest, "failed to parse volume_id to uuid"))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil && !errs.Is(err, http.ErrMissingFile) {
		errors.Handle(c, status.Error(code.BadRequest, "failed to get file"))
		return
	}

	size, file, err := h.openFile(fileHeader)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	entry, err := h.entryUC.Create(ctx, accountID, volumeID, req.Key, size, file)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusCreated, builder.ToEntryResponse(entry))
}

func (h *entryHandler) Update(c *gin.Context) {
	var req schema.UpdateEntryRequest
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

	entry, err := h.entryUC.Update(ctx, accountID, id, req.Key)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToEntryResponse(entry))
}

func (h *entryHandler) Delete(c *gin.Context) {
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

	if err := h.entryUC.Delete(ctx, accountID, id); err != nil {
		errors.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *entryHandler) openFile(fileHeader *multipart.FileHeader) (uint64, multipart.File, error) {
	if fileHeader == nil {
		return 0, nil, nil
	}

	if fileHeader.Size < 0 {
		return 0, nil, status.Error(code.BadRequest, "file is corrupted")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return 0, nil, status.Error(code.BadRequest, "failed to open file")
	}

	return uint64(fileHeader.Size), file, nil
}
