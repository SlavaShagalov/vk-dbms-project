package main

import (
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	"github.com/SlavaShagalov/vk-dbms-project/internal/pkg/db"
	pkgLog "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/log/zap"

	userDelivery "github.com/SlavaShagalov/vk-dbms-project/internal/user/delivery/http"
	userRepository "github.com/SlavaShagalov/vk-dbms-project/internal/user/repository/pgx"
	userService "github.com/SlavaShagalov/vk-dbms-project/internal/user/service"
)

func main() {
	// Logger
	logger, logfile, err := pkgLog.NewProdLogger()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Println(err)
		}
		err = logfile.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// Context
	//ctx := context.Background()

	//Database
	pool, err := db.NewPgxPool(logger)
	if err != nil {
		os.Exit(1)
	}
	defer pool.Close()

	// Repositories
	usersRepo := userRepository.NewRepository(pool, logger)

	// Services
	userServ := userService.NewService(usersRepo, logger)

	// Router
	router := httprouter.New()

	// Delivery
	userDelivery.RegisterHandlers(router, logger, userServ)

	// Server
	server := http.Server{
		Addr:    ":5000",
		Handler: router,
	}
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
