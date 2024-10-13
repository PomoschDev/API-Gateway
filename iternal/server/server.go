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

	//Эндпоинты usersPrivate
	usersPrivateRoute := route.r.PathPrefix(getEndpoint("users")).Subrouter()
	usersPrivateRoute.Use(cors.Default().Handler, route.authMiddleware)

	//Эндпоинты usersPrivate
	usersPublicRoute := route.r.PathPrefix(getEndpoint("users")).Subrouter()
	usersPublicRoute.Use(cors.Default().Handler, route.publicMiddleware)

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
		//Приватные
		{
			usersPrivateRoute.HandleFunc("", route.GetUsers).Methods(http.MethodGet, http.MethodOptions)
			usersPrivateRoute.HandleFunc("", route.CreateUser).Methods(http.MethodPost, http.MethodOptions)
			usersPrivateRoute.HandleFunc("/{id:[0-9]+}", route.ChangeUserType).Methods(http.MethodPatch, http.MethodOptions)
		}

		//Публичные
		{
			usersPublicRoute.HandleFunc("/{id:[0-9]+}", route.GetUser).Methods(http.MethodGet, http.MethodOptions)
			usersPublicRoute.HandleFunc("/{id:[0-9]+}", route.UpdateUser).Methods(http.MethodPut, http.MethodOptions)
			usersPublicRoute.HandleFunc("/{id:[0-9]+}", route.DeleteUserByID).Methods(http.MethodDelete, http.MethodOptions)
			usersPublicRoute.HandleFunc("/isExists", route.UserIsExists).Methods(http.MethodPost, http.MethodOptions)
			usersPublicRoute.HandleFunc("/isRole", route.UserIsRole).Methods(http.MethodPost, http.MethodOptions)
			usersPublicRoute.HandleFunc("/comparePassword", route.ComparePassword).Methods(http.MethodPost, http.MethodOptions)
			usersPublicRoute.HandleFunc("/", route.FindUserByEmail).Queries("email", "{email}").Methods(http.MethodGet,
				http.MethodOptions)
			usersPublicRoute.HandleFunc("/", route.FindUserByPhone).Queries("phone", "{phone}").Methods(http.MethodGet,
				http.MethodOptions)
		}
	}

	route.r.Use(cors.Default().Handler, mux.CORSMethodMiddleware(route.r))

	// CORS обработчик
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
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
