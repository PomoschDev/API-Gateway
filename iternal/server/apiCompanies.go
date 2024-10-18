package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/utilities"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Companies godoc
// @Summary      Список всех компаний
// @Description  Массив компаний в базе данных
// @Tags         Company
// @Accept       json
// @Produce      json
// @Success      200  {object}  DatabaseServicev1.CompaniesResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/companies [get]
func (route Router) Companies(w http.ResponseWriter, r *http.Request) {
	response, err := route.databaseService.Companies(r.Context(), nil)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// CreateCompany godoc
// @Summary      Создание компании
// @Description  Создание новой сущности компании
// @Tags         Company
// @Accept       json
// @Produce      json
// @Param        company body DatabaseServicev1.CreateCompanyRequest false "Сущность компании"
// @Success      200  {object}  DatabaseServicev1.Company
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/companies [post]
func (route Router) CreateCompany(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.CreateCompanyRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	user, err := route.databaseService.FindUserById(r.Context(), &DatabaseServicev1.FindUserByIdRequest{Id: request.UserId})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	if user.Type == 0 {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetHTTPError(w, "У физических лиц не может быть компании", http.StatusNotFound)
		return
	}

	response, err := route.databaseService.CreateCompany(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// Company godoc
// @Summary      Поиск компании
// @Description  Поиск компании по ID
// @Tags         Company
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Company ID"
// @Success      200  {object}  DatabaseServicev1.Company
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/companies/{id} [get]
func (route Router) Company(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindCompanyByIdRequest{Id: id}

	response, err := route.databaseService.FindCompanyById(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// FindCompanyByPhone godoc
// @Summary      Поиск компании по номеру телефона
// @Description  Поиск компании по ее phone
// @Tags         Company
// @Accept       json
// @Produce      json
// @Param        phone query string true "Phone"
// @Success      200  {object}  DatabaseServicev1.Company
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/companies/ [get]
func (route Router) FindCompanyByPhone(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")

	if phone == "" || len(phone) == 0 {
		logger.Error("Phone не может быть пустым")
		SetHTTPError(w, "Поле \"phone\" не может быть пустым", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindCompanyByIdPhoneRequest{Phone: r.URL.Query().Get("phone")}

	logger.Info("request: %+v", request)

	response, err := route.databaseService.FindCompanyByPhone(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// FindCompanyCard godoc
// @Summary      Извлечение банковской карты компании
// @Description  Извлечение банковской карты компании по ее ID
// @Tags         Company
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Company ID"
// @Success      200  {object}  DatabaseServicev1.CardCompany
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/companies/{id}/card [get]
func (route Router) FindCompanyCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindCompanyCardRequest{Id: id}

	response, err := route.databaseService.FindCompanyCard(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// DeleteCompanyByModel godoc
// @Summary      Удаление компании по модели
// @Description  Удаляет компании опираясь на всю сущность модели
// @Tags         Company
// @Accept       json
// @Produce      json
// @Param        company body DatabaseServicev1.DeleteCompanyByModelRequest false "Модель компании"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/companies/deleteModel [post]
func (route Router) DeleteCompanyByModel(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.DeleteCompanyByModelRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.DeleteCompanyByModel(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// DeleteCompanyByID godoc
// @Summary      Удаление компании
// @Description  Удаление компании по ID
// @Tags         Company
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID компании"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/company/{id} [delete]
func (route Router) DeleteCompanyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	request := &DatabaseServicev1.DeleteCompanyByIdRequest{Id: id}

	response, err := route.databaseService.DeleteCompanyById(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// UpdateCompany godoc
// @Summary      Обновление компании
// @Description  Обновление сущности компании
// @Tags         Company
// @Accept       json
// @Produce      json
// @Param        company body DatabaseServicev1.UpdateUserRequest true "Модель для обновления"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id} [put]
func (route Router) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	id := utilities.StrToUint(mux.Vars(r)["id"])

	request := new(DatabaseServicev1.UpdateCompanyRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	request.Company.Id = id

	response, err := route.databaseService.UpdateCompany(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// AddCardToCompany godoc
// @Summary      Добавляет банковскую карту компании
// @Description  Добавляет банковскую карту компании, поле companyId это ID компании в базе данных, которой будем добавлять карту
// @Tags         Company
// @Accept       json
// @Produce      json
// @Param        card body DatabaseServicev1.AddCardToCompanyRequest false "Сущность банковской карты"
// @Success      200  {object}  DatabaseServicev1.AddCardToCompanyResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/companies/addCard [post]
func (route Router) AddCardToCompany(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.AddCardToCompanyRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.AddCardToCompany(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}
