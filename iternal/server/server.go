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

	//Эндпоинты usersPublic
	usersPublicRoute := route.r.PathPrefix(getEndpoint("users")).Subrouter()
	usersPublicRoute.Use(cors.Default().Handler, route.publicMiddleware)

	//Эндпоинты companiesPrivate
	companiesPrivateRoute := route.r.PathPrefix(getEndpoint("companies")).Subrouter()
	companiesPrivateRoute.Use(cors.Default().Handler, route.authMiddleware)

	//Эндпоинты companiesPublic
	companiesPublicRoute := route.r.PathPrefix(getEndpoint("companies")).Subrouter()
	companiesPublicRoute.Use(cors.Default().Handler, route.publicMiddleware)

	//Эндпоинты cardsPrivate
	cardsPrivateRoute := route.r.PathPrefix(getEndpoint("cards")).Subrouter()
	cardsPrivateRoute.Use(cors.Default().Handler, route.authMiddleware)

	//Эндпоинты companiesPublic
	cardsPublicRoute := route.r.PathPrefix(getEndpoint("cards")).Subrouter()
	cardsPublicRoute.Use(cors.Default().Handler, route.publicMiddleware)

	//Эндпоинты cardCompaniesPrivate
	cardCompaniesPrivateRoute := route.r.PathPrefix(getEndpoint("card/company")).Subrouter()
	cardCompaniesPrivateRoute.Use(cors.Default().Handler, route.authMiddleware)

	//Эндпоинты cardCompaniesPublic
	cardCompaniesPublicRoute := route.r.PathPrefix(getEndpoint("card/company")).Subrouter()
	cardCompaniesPublicRoute.Use(cors.Default().Handler, route.publicMiddleware)

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
			usersPrivateRoute.HandleFunc("/{id:[0-9]+}/company", route.FindUserCompany).Methods(http.MethodGet,
				http.MethodOptions)
			usersPrivateRoute.HandleFunc("/{id:[0-9]+}/donation", route.FindUserDonations).Methods(http.MethodGet,
				http.MethodOptions)
			usersPrivateRoute.HandleFunc("/{id:[0-9]+}/card", route.FindUserCard).Methods(http.MethodGet,
				http.MethodOptions)
			usersPrivateRoute.HandleFunc("/addCard", route.AddCardToUser).Methods(http.MethodPost,
				http.MethodOptions)
			/*usersPrivateRoute.HandleFunc("/addCard", route.AddCardToUser).Methods(http.MethodPost,
			http.MethodOptions)*/
			usersPrivateRoute.HandleFunc("/deleteModel", route.DeleteUserByModel).Methods(http.MethodPost,
				http.MethodOptions)
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

	//Компании
	{
		//Приватные
		{
			companiesPrivateRoute.HandleFunc("", route.Companies).Methods(http.MethodGet, http.MethodOptions)
			companiesPrivateRoute.HandleFunc("/{id:[0-9]+}", route.Company).Methods(http.MethodGet, http.MethodOptions)
			companiesPrivateRoute.HandleFunc("/{id:[0-9]+}/card", route.FindCompanyCard).Methods(http.MethodGet,
				http.MethodOptions)
			companiesPrivateRoute.HandleFunc("/deleteModel", route.DeleteCompanyByModel).Methods(http.MethodPost,
				http.MethodOptions)
			companiesPrivateRoute.HandleFunc("/{id:[0-9]+}", route.DeleteCompanyByID).Methods(http.MethodDelete, http.MethodOptions)
			companiesPrivateRoute.HandleFunc("/{id:[0-9]+}", route.UpdateCompany).Methods(http.MethodPut, http.MethodOptions)
			companiesPrivateRoute.HandleFunc("/addCard", route.AddCardToUser).Methods(http.MethodPost,
				http.MethodOptions)
		}

		//Публичные
		{
			companiesPublicRoute.HandleFunc("", route.CreateCompany).Methods(http.MethodPost, http.MethodOptions)
			companiesPublicRoute.HandleFunc("/", route.FindCompanyByPhone).Queries("phone", "{phone}").Methods(http.MethodGet,
				http.MethodOptions)
		}
	}

	//Банковские карты пользователей
	{
		//Приватные
		{
			cardsPrivateRoute.HandleFunc("", route.Cards).Methods(http.MethodGet, http.MethodOptions)
			cardsPrivateRoute.HandleFunc("/{id:[0-9]+}", route.Card).Methods(http.MethodGet, http.MethodOptions)
			cardsPrivateRoute.HandleFunc("/deleteModel", route.DeleteCardByModel).Methods(http.MethodPost,
				http.MethodOptions)
			cardsPrivateRoute.HandleFunc("/{id:[0-9]+}", route.DeleteCardById).Methods(http.MethodDelete,
				http.MethodOptions)
			cardsPrivateRoute.HandleFunc("/{id:[0-9]+}", route.UpdateCard).Methods(http.MethodPut, http.MethodOptions)
		}

		//Публичные
		{
			cardsPublicRoute.HandleFunc("", route.CreateCard).Methods(http.MethodPost, http.MethodOptions)
		}
	}

	//Банковские карты компаний
	{
		//Приватные
		{
			cardCompaniesPrivateRoute.HandleFunc("", route.CardCompanies).Methods(http.MethodGet, http.MethodOptions)
			cardCompaniesPrivateRoute.HandleFunc("/{id:[0-9]+}", route.CardCompany).Methods(http.MethodGet,
				http.MethodOptions)
			cardCompaniesPrivateRoute.HandleFunc("/deleteModel", route.DeleteCardCompaniesByModel).Methods(http.
				MethodPost,
				http.MethodOptions)
			cardCompaniesPrivateRoute.HandleFunc("/{id:[0-9]+}", route.DeleteCardCompanyById).Methods(http.MethodDelete,
				http.MethodOptions)
			cardCompaniesPrivateRoute.HandleFunc("/{id:[0-9]+}", route.UpdateCardCompany).Methods(http.MethodPut,
				http.MethodOptions)
		}

		//Публичные
		{
			cardCompaniesPublicRoute.HandleFunc("", route.CreateCardCompany).Methods(http.MethodPost,
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
