package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/utilities"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Donations godoc
// @Summary      Список всех пожертвований
// @Description  Список всех пожертвований в базе данных
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Success      200  {object}  DatabaseServicev1.DonationsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/donations [get]
func (route Router) Donations(w http.ResponseWriter, r *http.Request) {
	response, err := route.databaseService.Donations(r.Context(), nil)
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

// CreateDonation godoc
// @Summary      Создание пожертвования
// @Description  Создание пожертвования
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Param        donation body DatabaseServicev1.CreateDonationsRequest false "Сущность пожертвования"
// @Success      200  {object}  DatabaseServicev1.Card
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/donations [post]
func (route Router) CreateDonation(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.CreateDonationsRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	if request.GetUserId() <= 0 {
		SetHTTPError(w, "Поле \"UserID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	if request.GetWardId() <= 0 {
		SetHTTPError(w, "Поле \"WardID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	if request.GetAmount() <= 0 {
		SetHTTPError(w, "Поле \"Amount\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	_, err := route.databaseService.FindUserById(r.Context(), &DatabaseServicev1.FindUserByIdRequest{Id: request.GetUserId()})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	response, err := route.databaseService.CreateDonations(r.Context(), request)
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

// FindDonationWards godoc
// @Summary      Извлечение подопечных пожертвования
// @Description  Извлечение подопечных пожертвования по ID пожертвования
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Donation ID"
// @Success      200  {object}  DatabaseServicev1.FindDonationWardsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/donations/{id}/wards [get]
func (route Router) FindDonationWards(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindDonationWardsRequest{Id: id}

	response, err := route.databaseService.FindDonationWards(r.Context(), request)
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

// FindDonationUser godoc
// @Summary      Извлечение пользователя из пожертвования
// @Description  Извлечение пользователя из пожертвования по ID пожертвования
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Donation ID"
// @Success      200  {object}  DatabaseServicev1.FindDonationWardsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/donations/{id}/user [get]
func (route Router) FindDonationUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindDonationUserRequest{Id: id}

	response, err := route.databaseService.FindDonationUser(r.Context(), request)
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

// Donation godoc
// @Summary      Поиск пожертвования
// @Description  Поиск пожертвовани по ID
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Donation ID"
// @Success      200  {object}  DatabaseServicev1.CreateDonationsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/donations/{id} [get]
func (route Router) Donation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindDonationByIdRequest{Id: id}

	response, err := route.databaseService.FindDonationById(r.Context(), request)
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

// DeleteDonationByModel godoc
// @Summary      Удаление пожертвования по модели
// @Description  Удаляет пожертвование опираясь на всю сущность модели
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        donation body DatabaseServicev1.DeleteDonationByModelRequest false "Сущность пожертвования"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/donations/deleteModel [post]
func (route Router) DeleteDonationByModel(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.DeleteDonationByModelRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.DeleteDonationByModel(r.Context(), request)
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

// DeleteDonationById godoc
// @Summary      Удаление пожертвования
// @Description  Удаление пожертвования по ID
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID пожертвования"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/donations/{id} [delete]
func (route Router) DeleteDonationById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.DeleteDonationByIdRequest{Id: id}

	response, err := route.databaseService.DeleteDonationById(r.Context(), request)
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

// UpdateDonation godoc
// @Summary      Обновление пожертвования
// @Description  Обновление пожертвования
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        donation body DatabaseServicev1.UpdateDonationsRequest false "Модель для обновления"
// @Success      200  {object}  DatabaseServicev1.CreateDonationsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/donations [put]
func (route Router) UpdateDonation(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.UpdateDonationsRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.UpdateDonation(r.Context(), request)
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
