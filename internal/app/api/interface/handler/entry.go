package handler

import (
	errs "errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

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
	Head(*gin.Context)
	Get(*gin.Context)
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

	entry, err := h.entryUC.Create(ctx, accountID, req.VolumeName, req.Key, size, file)
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

	volumeName := c.Param("volumeName")
	key := strings.TrimPrefix(c.Param("key"), "/")

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	entry, err := h.entryUC.Update(ctx, accountID, volumeName, key, req.Key)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToEntryResponse(entry))
}

func (h *entryHandler) Delete(c *gin.Context) {
	volumeName := c.Param("volumeName")
	key := strings.TrimPrefix(c.Param("key"), "/")

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	if err := h.entryUC.Delete(ctx, accountID, volumeName, key); err != nil {
		errors.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *entryHandler) Head(c *gin.Context) {
	volumeName := c.Param("volumeName")
	key := strings.TrimPrefix(c.Param("key"), "/")

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		log.Println(err)
		c.Header("Content-Length", "0")
		c.Status(errors.GetStatusCode(err))
		return
	}

	ctx := c.Request.Context()

	entry, err := h.entryUC.Head(ctx, accountID, volumeName, key)
	if err != nil {
		log.Println(err)
		c.Header("Content-Length", "0")
		c.Status(errors.GetStatusCode(err))
		return
	}

	c.Header("Content-Length", strconv.FormatUint(entry.Size, 10))
	c.Header("Content-Type", entry.Type)
	c.Header("Last-Modified", entry.UpdatedAt.Format(http.TimeFormat))

	c.Status(http.StatusOK)
}

func (h *entryHandler) Get(c *gin.Context) {
	volumeName := c.Param("volumeName")
	key := strings.TrimPrefix(c.Param("key"), "/")

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	entry, body, err := h.entryUC.GetOne(ctx, accountID, volumeName, key)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if body == nil {
		c.Header("Content-Length", strconv.FormatUint(entry.Size, 10))
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Last-Modified", entry.UpdatedAt.Format(http.TimeFormat))
		return
	}

	defer func() {
		if err := body.Close(); err != nil {
			errors.Handle(c, err)
			return
		}
	}()

	c.Header("Content-Length", strconv.FormatUint(entry.Size, 10))
	c.Header("Content-Type", entry.Type)
	c.Header("Last-Modified", entry.UpdatedAt.Format(http.TimeFormat))

	if _, err := io.Copy(c.Writer, body); err != nil {
		errors.Handle(c, err)
		return
	}
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
