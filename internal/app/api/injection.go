package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/file"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/handler"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/middleware"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
)

var (
	authorizationMW middleware.AuthorizationMiddleware

	healthHdl handler.HealthHandler
	volumeHdl handler.VolumeHandler
	entryHdl  handler.EntryHandler
)

func inject(db *sqlx.DB, fs afero.Fs, config *serverConfig) {
	transactionObj := transaction.NewDBTransactionObject(db)

	accountRepo := api.NewAccountRepository(&http.Client{}, "http://account-api:8000/authorization")
	volumeRepo := database.NewVolumeRepository(db)
	entryRepo := database.NewEntryRepository(db)
	bodyRepo := file.NewBodyRepository(fs, config.fileSystem.BasePath)

	volumeServ := service.NewVolumeService(volumeRepo)
	entryServ := service.NewEntryService(entryRepo)

	authorizationUC := usecase.NewAuthorizationUsecase(accountRepo)
	volumeUC := usecase.NewVolumeUsecase(transactionObj, volumeRepo, volumeServ)
	entryUC := usecase.NewEntryUsecase(transactionObj, entryRepo, bodyRepo, volumeRepo, entryServ)

	authorizationMW = middleware.NewAuthorizationMiddleware(authorizationUC)

	healthHdl = handler.NewHealthHandler()
	volumeHdl = handler.NewVolumeHandler(volumeUC)
	entryHdl = handler.NewEntryHandler(entryUC)
}
