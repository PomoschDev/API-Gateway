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
	cfg             *config.Config
}

const apiStr = "/api/v1/"

// New - создает новый роутер для маршрутизации
func New(cfg *config.Config, grpcClient *grpc.Api) *http.Server {
	router := &Router{
		r:               mux.NewRouter(),
		mu:              sync.Mutex{},
		databaseService: grpcClient.Client,
		cfg:             cfg,
	}

	return router.loadEndpoints()
}

func getEndpoint(endpoint string) string {
	return fmt.Sprintf("%s%s", apiStr, endpoint)
}

func (route *Router) loadEndpoints() *http.Server {
	addr := fmt.Sprintf(":%d", route.cfg.APIServer.Port)

	//Эндпоинты auth
	authRoute := route.r.PathPrefix(getEndpoint("auth")).Subrouter()

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

	//Эндпоинты donationsPrivate
	donationsPrivateRoute := route.r.PathPrefix(getEndpoint("donations")).Subrouter()
	donationsPrivateRoute.Use(cors.Default().Handler, route.authMiddleware)

	//Эндпоинты donationsPublic
	donationsPublicRoute := route.r.PathPrefix(getEndpoint("donations")).Subrouter()
	donationsPublicRoute.Use(cors.Default().Handler, route.publicMiddleware)

	//Эндпоинты wardsPrivate
	wardsPrivateRoute := route.r.PathPrefix(getEndpoint("wards")).Subrouter()
	wardsPrivateRoute.Use(cors.Default().Handler, route.authMiddleware)

	//Эндпоинты wardsPublic
	wardsPublicRoute := route.r.PathPrefix(getEndpoint("wards")).Subrouter()
	wardsPublicRoute.Use(cors.Default().Handler, route.publicMiddleware)

	//Эндпоинты wardsPrivate
	paymentPrivateRoute := route.r.PathPrefix(getEndpoint("payment")).Subrouter()
	paymentPrivateRoute.Use(cors.Default().Handler, route.authMiddleware)

	//Эндпоинты wardsPublic
	paymentPublicRoute := route.r.PathPrefix(getEndpoint("payment")).Subrouter()
	paymentPublicRoute.Use(cors.Default().Handler, route.publicMiddleware)

	//Swagger
	{
		if route.cfg.Swagger {
			route.r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
				httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
				httpSwagger.DeepLinking(true),
				httpSwagger.DocExpansion("none"),
				httpSwagger.DomID("swagger-ui"),
			)).Methods(http.MethodGet)
		}
	}

	//Аутентификация
	{
		authRoute.HandleFunc("/login", route.Login).Methods(http.MethodPost, http.MethodOptions)
		authRoute.HandleFunc("/registration", route.Registration).Methods(http.MethodPost, http.MethodOptions)
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
			usersPrivateRoute.HandleFunc("/deleteModel", route.DeleteUserByModel).Methods(http.MethodPost,
				http.MethodOptions)
			usersPrivateRoute.HandleFunc("/{id:[0-9]+}/photo", route.DeleteUserPhoto).Methods(http.MethodDelete,
				http.MethodOptions)
			usersPrivateRoute.HandleFunc("/{id:[0-9]+}/photo", route.SetUserPhoto).Methods(http.MethodPost,
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
			usersPublicRoute.HandleFunc("/{id:[0-9]+}/photo", route.GetUserPhoto).Methods(http.MethodGet, http.MethodOptions)
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
			companiesPrivateRoute.HandleFunc("", route.UpdateCompany).Methods(http.MethodPut, http.MethodOptions)
			companiesPrivateRoute.HandleFunc("/addCard", route.AddCardToCompany).Methods(http.MethodPost,
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

	//Пожертвования
	{
		//Приватные
		{
			donationsPrivateRoute.HandleFunc("/{id:[0-9]+}/wards", route.FindDonationWards).Methods(http.MethodGet,
				http.MethodOptions)
			donationsPrivateRoute.HandleFunc("/{id:[0-9]+}/user", route.FindDonationUser).Methods(http.MethodGet,
				http.MethodOptions)
			donationsPrivateRoute.HandleFunc("/{id:[0-9]+}", route.Donation).Methods(http.MethodGet,
				http.MethodOptions)
			donationsPrivateRoute.HandleFunc("/deleteModel", route.DeleteDonationByModel).Methods(http.
				MethodPost,
				http.MethodOptions)
			donationsPrivateRoute.HandleFunc("/{id:[0-9]+}", route.DeleteDonationById).Methods(http.MethodDelete,
				http.MethodOptions)
			donationsPrivateRoute.HandleFunc("", route.UpdateDonation).Methods(http.MethodPut,
				http.MethodOptions)
		}

		//Публичные
		{
			donationsPublicRoute.HandleFunc("", route.CreateDonation).Methods(http.MethodPost,
				http.MethodOptions)
			donationsPublicRoute.HandleFunc("", route.Donations).Methods(http.MethodGet, http.MethodOptions)
		}
	}

	//Подопечные
	{
		//Приватные
		{
			wardsPrivateRoute.HandleFunc("", route.CreateWard).Methods(http.MethodPost,
				http.MethodOptions)
			wardsPrivateRoute.HandleFunc("/deleteModel", route.DeleteWardByModel).Methods(http.
				MethodPost,
				http.MethodOptions)
			wardsPrivateRoute.HandleFunc("/{id:[0-9]+}", route.DeleteWardById).Methods(http.MethodDelete,
				http.MethodOptions)
			wardsPrivateRoute.HandleFunc("", route.UpdateWard).Methods(http.MethodPut,
				http.MethodOptions)
			wardsPrivateRoute.HandleFunc("/{id:[0-9]+}/donations", route.FindWardDonations).Methods(http.MethodGet,
				http.MethodOptions)
		}

		//Публичные
		{
			wardsPublicRoute.HandleFunc("", route.Wards).Methods(http.MethodGet, http.MethodOptions)
			wardsPublicRoute.HandleFunc("/{id:[0-9]+}", route.Ward).Methods(http.MethodGet,
				http.MethodOptions)
		}
	}

	//Платежи
	{
		//Приватные
		{
			paymentPrivateRoute.HandleFunc("", route.Payment).Methods(http.MethodPost, http.MethodOptions)
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
		WriteTimeout: route.cfg.APIServer.Timeout,
		ReadTimeout:  route.cfg.APIServer.Timeout,
		IdleTimeout:  route.cfg.APIServer.Timeout,
		Handler:      cors.AllowAll().Handler(handler),
	}

	return srv
}
