package main

import (
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
	"os"

	"github.com/SlavaShagalov/vk-dbms-project/internal/pkg/db"
	pkgLog "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/log/zap"

	serviceDelivery "github.com/SlavaShagalov/vk-dbms-project/internal/service/delivery/http"
	serviceRepository "github.com/SlavaShagalov/vk-dbms-project/internal/service/repository/pgx"
	serviceService "github.com/SlavaShagalov/vk-dbms-project/internal/service/service"

	postDelivery "github.com/SlavaShagalov/vk-dbms-project/internal/post/delivery/http"
	postRepository "github.com/SlavaShagalov/vk-dbms-project/internal/post/repository/pgx"
	postService "github.com/SlavaShagalov/vk-dbms-project/internal/post/service"

	threadDelivery "github.com/SlavaShagalov/vk-dbms-project/internal/thread/delivery/http"
	threadRepository "github.com/SlavaShagalov/vk-dbms-project/internal/thread/repository/pgx"
	threadService "github.com/SlavaShagalov/vk-dbms-project/internal/thread/service"

	forumDelivery "github.com/SlavaShagalov/vk-dbms-project/internal/forum/delivery/http"
	forumRepository "github.com/SlavaShagalov/vk-dbms-project/internal/forum/repository/pgx"
	forumService "github.com/SlavaShagalov/vk-dbms-project/internal/forum/service"

	userDelivery "github.com/SlavaShagalov/vk-dbms-project/internal/user/delivery/http"
	userRepository "github.com/SlavaShagalov/vk-dbms-project/internal/user/repository/pgx"
	userService "github.com/SlavaShagalov/vk-dbms-project/internal/user/service"
)

func main() {
	// Logger
	logger := pkgLog.NewProdLogger()

	//Database
	pool, err := db.NewPgxPool(logger)
	if err != nil {
		os.Exit(1)
	}
	defer pool.Close()

	// Repositories
	userRepo := userRepository.NewRepository(pool, logger)
	forumRepo := forumRepository.NewRepository(pool, logger)
	threadRepo := threadRepository.NewRepository(pool, logger)
	postRepo := postRepository.NewRepository(pool, logger)
	serviceRepo := serviceRepository.NewRepository(pool, logger)

	// Services
	userServ := userService.NewService(userRepo, logger)
	forumServ := forumService.NewService(forumRepo, logger)
	threadServ := threadService.NewService(threadRepo, logger)
	postServ := postService.NewService(postRepo, logger)
	serviceServ := serviceService.NewService(serviceRepo, logger)

	// Router
	router := httprouter.New()

	// Delivery
	userDelivery.RegisterHandlers(router, logger, userServ)
	threadDelivery.RegisterHandlers(router, logger, threadServ)
	forumDelivery.RegisterHandlers(router, logger, forumServ)
	postDelivery.RegisterHandlers(router, logger, postServ)
	serviceDelivery.RegisterHandlers(router, logger, serviceServ)

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
