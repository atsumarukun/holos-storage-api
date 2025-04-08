package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Serve() {
	conf := loadServerConfig()

	db, err := NewDatabase(&conf.database)
	if err != nil {
		log.Fatalln(err.Error())
	}

	inject(db)

	r := gin.Default()
	registerRouter(r)

	srv := &http.Server{
		Addr:              ":8000",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	<-ctx.Done()

	ctx, stop = context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err.Error())
	}
}
