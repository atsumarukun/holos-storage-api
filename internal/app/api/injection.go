package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/handler"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/middleware"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
)

var (
	authorizationMW middleware.AuthorizationMiddleware

	healthHdl handler.HealthHandler
	volumeHdl handler.VolumeHandler
)

func inject(db *sqlx.DB) {
	transactionObj := transaction.NewDBTransactionObject(db)

	accountRepo := api.NewAccountRepository(&http.Client{}, "http://account-api:8000/authorization")
	volumeRepo := database.NewVolumeRepository(db)

	volumeServ := service.NewVolumeService(volumeRepo)

	authorizationUC := usecase.NewAuthorizationUsecase(accountRepo)
	volumeUC := usecase.NewVolumeUsecase(transactionObj, volumeRepo, volumeServ)

	authorizationMW = middleware.NewAuthorizationMiddleware(authorizationUC)

	healthHdl = handler.NewHealthHandler()
	volumeHdl = handler.NewVolumeHandler(volumeUC)
}
