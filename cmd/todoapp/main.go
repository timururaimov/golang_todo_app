package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_postgers_pool "github.com/timururaimov/golang_todo_app/internal/core/repository/postgres/conn"
	core_http_middleware "github.com/timururaimov/golang_todo_app/internal/core/transport/http/middleware"
	core_http_server "github.com/timururaimov/golang_todo_app/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/timururaimov/golang_todo_app/internal/features/tasks/repository/postgres"
	tasks_service "github.com/timururaimov/golang_todo_app/internal/features/tasks/service"
	tasks_transport_http "github.com/timururaimov/golang_todo_app/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/timururaimov/golang_todo_app/internal/features/users/repository/postgres"
	users_service "github.com/timururaimov/golang_todo_app/internal/features/users/service"
	users_transport_http "github.com/timururaimov/golang_todo_app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)

	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init aplecation logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	pool, err := core_postgers_pool.NewConnectionPool(
		ctx,
		core_postgers_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("faild to connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("aplicatation timeZone", zap.Any("zone", timeZone))
	//Users
	logger.Debug("intializing feature", zap.String("featur", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	//Tasks
	logger.Debug("intializing feature", zap.String("featur", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksSrevice := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHttp := tasks_transport_http.NewTasksHTTPHandler(tasksSrevice)

	httpSever := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.ReqisterRoutes(usersTransportHttp.Routes()...)
	apiVersionRouterV1.ReqisterRoutes(tasksTransportHttp.Routes()...)

	// apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
	// 	core_http_server.ApiVersion2,
	// 	какой новый midlware в новый верссиии api
	// 	core_http_middleware.Panic(),
	// )
	// apiVersionRouterV2.ReqisterRoutes(usersTransportHttp.Routes()...)

	httpSever.RegisterAPIRouters(
		apiVersionRouterV1,
		// apiVersionRouterV2
	)

	if err := httpSever.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}

}
