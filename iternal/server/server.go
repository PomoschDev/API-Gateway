package server

import (
	_ "apiGateway/docs"
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/iternal/grpc"
	"apiGateway/pkg/config"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"sync"
)

// Router - сущность маршрутизатора, содержит приватные поля для работы исключительно внутри пакета
type Router struct {
	r               *mux.Router
	mu              sync.Mutex
	databaseService DatabaseServicev1.DatabaseServiceClient
}

const apiStr = "/api/v1/"

// New - создает новый роутер для маршрутизации
func New(cfg *config.Config, grpcClient *grpc.Api) *http.Server {
	router := &Router{
		r:               mux.NewRouter(),
		mu:              sync.Mutex{},
		databaseService: grpcClient.Client,
	}

	return router.loadEndpoints(cfg)
}

func getEndpoint(endpoint string) string {
	return fmt.Sprintf("%s%s", apiStr, endpoint)
}

func (route *Router) loadEndpoints(cfg *config.Config) *http.Server {
	addr := fmt.Sprintf(":%d", cfg.APIServer.Port)

	privateRouter := route.r.PathPrefix(apiStr).Subrouter()
	privateRouter.Use(cors.Default().Handler, route.middleware)

	usersRoute := route.r.PathPrefix(getEndpoint("users")).Subrouter()
	usersRoute.Use(cors.Default().Handler, route.middleware)

	//Swagger
	{
		route.r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)).Methods(http.MethodGet)
	}

	//Пользователи
	{
		usersRoute.HandleFunc("", route.GetUsers).Methods(http.MethodGet, http.MethodOptions)
		usersRoute.HandleFunc("/{id:[0-9]+}", route.GetUser).Methods(http.MethodGet, http.MethodOptions)
		usersRoute.HandleFunc("/{id:[0-9]+}", route.UpdateUser).Methods(http.MethodPut, http.MethodOptions)
		usersRoute.HandleFunc("", route.CreateUser).Methods(http.MethodPost, http.MethodOptions)
		usersRoute.HandleFunc("/{id:[0-9]+}", route.DeleteUserByID).Methods(http.MethodDelete, http.MethodOptions)

	}

	route.r.Use(cors.Default().Handler, mux.CORSMethodMiddleware(route.r))

	// CORS обработчик
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "application/json"},
	})
	handler := crs.Handler(route.r)

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: cfg.APIServer.Timeout,
		ReadTimeout:  cfg.APIServer.Timeout,
		IdleTimeout:  cfg.APIServer.Timeout,
		Handler:      cors.AllowAll().Handler(handler),
	}

	return srv
}
