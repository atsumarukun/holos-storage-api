package handler

import (
	errs "errors"
	"io"
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
	var size uint64
	var body io.Reader
	volumeID, key, fileHeader, err := h.parseCreateRequest(c)
	if err == nil {
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
		if fileHeader.Size < 0 {
			errors.Handle(c, status.Error(code.BadRequest, "invalid file"))
			return
		}
		size = uint64(fileHeader.Size)
		body = file
	} else if !errs.Is(err, http.ErrMissingFile) {
		errors.Handle(c, err)
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	entry, err := h.entryUC.Create(ctx, accountID, volumeID, key, size, body)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusCreated, builder.ToEntryResponse(entry))
}

func (h *entryHandler) parseCreateRequest(c *gin.Context) (volumeID uuid.UUID, key string, fileHeader *multipart.FileHeader, err error) {
	var req schema.CreateEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		return uuid.Nil, "", nil, status.Error(code.BadRequest, "failed to parse multipart/form-data")
	}

	volumeID, err = uuid.Parse(req.VolumeID)
	if err != nil {
		return uuid.Nil, "", nil, status.Error(code.BadRequest, "failed to parse volume_id to uuid")
	}

	fileHeader, err = c.FormFile("file")
	return volumeID, req.Key, fileHeader, err
}
