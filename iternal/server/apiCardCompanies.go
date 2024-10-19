package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/utilities"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// CardCompanies godoc
// @Summary      Банковская карта компании
// @Description  Банковская карта компании в базе данных
// @Tags         CardCompany
// @Accept       json
// @Produce      json
// @Success      200  {object}  DatabaseServicev1.CardsCompaniesResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/card/company [get]
func (route Router) CardCompanies(w http.ResponseWriter, r *http.Request) {
	response, err := route.databaseService.CardsCompanies(r.Context(), nil)
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

// CreateCardCompany godoc
// @Summary      Создание банковской карты компании
// @Description  Создание банковской карты компании
// @Tags         CardCompany
// @Accept       json
// @Produce      json
// @Param        card body DatabaseServicev1.CreateCardCompanyRequest false "Сущность банковской карты компании"
// @Success      200  {object}  DatabaseServicev1.Card
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/card/company [post]
func (route Router) CreateCardCompany(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.CreateCardCompanyRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	if request.GetCompanyId() <= 0 {
		SetHTTPError(w, "Поле \"UserID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	_, err := route.databaseService.FindCompanyById(r.Context(), &DatabaseServicev1.FindCompanyByIdRequest{Id: request.
		GetCompanyId()})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	response, err := route.databaseService.CreateCardCompany(r.Context(), request)
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

// CardCompany godoc
// @Summary      Поиск банковской карты компании
// @Description  Поиск банковской карты компании по ID
// @Tags         CardCompany
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Card ID"
// @Success      200  {object}  DatabaseServicev1.CardCompany
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/card/company/{id} [get]
func (route Router) CardCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindCardCompanyByIDRequest{Id: id}

	response, err := route.databaseService.FindCardCompanyByID(r.Context(), request)
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

// DeleteCardCompaniesByModel godoc
// @Summary      Удаление банковской карты компании по модели
// @Description  Удаляет банковскую карту компании опираясь на всю сущность модели
// @Tags         CardCompany
// @Accept       json
// @Produce      json
// @Param        card body DatabaseServicev1.CardCompany false "Сущность банковской карты компании"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/card/company/deleteModel [post]
func (route Router) DeleteCardCompaniesByModel(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.CardCompany)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.DeleteCardCompanyByModel(r.Context(), request)
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

// DeleteCardCompanyById godoc
// @Summary      Удаление банковской карты компании
// @Description  Удаление банковской карты компании по ID
// @Tags         CardCompany
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID банковской карты компании"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/card/company/{id} [delete]
func (route Router) DeleteCardCompanyById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.DeleteCardCompanyByIdRequest{Id: id}

	response, err := route.databaseService.DeleteCardCompanyById(r.Context(), request)
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

// UpdateCardCompany godoc
// @Summary      Обновление банковской карты компании
// @Description  Обновление банковской карты компании
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        card body DatabaseServicev1.CardCompany false "Модель для обновления"
// @Success      200  {object}  DatabaseServicev1.CardCompany
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/card/company [put]
func (route Router) UpdateCardCompany(w http.ResponseWriter, r *http.Request) {
	id := utilities.StrToUint(mux.Vars(r)["id"])

	request := new(DatabaseServicev1.CardCompany)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	request.Id = id

	response, err := route.databaseService.UpdateCardCompany(r.Context(), request)
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
