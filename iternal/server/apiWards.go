package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/utilities"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Wards godoc
// @Summary      Список всех подопечных
// @Description  Список всех подопечных в базе данных
// @Tags         Wards
// @Accept       json
// @Produce      json
// @Success      200  {object}  DatabaseServicev1.WardsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/wards [get]
func (route Router) Wards(w http.ResponseWriter, r *http.Request) {
	response, err := route.databaseService.Wards(r.Context(), nil)
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

// CreateWard godoc
// @Summary      Создание подопечного
// @Description  Создание подопечного
// @Tags         Wards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ward body DatabaseServicev1.CreateWardRequest false "Сущность подопечного"
// @Success      200  {object}  DatabaseServicev1.Ward
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/wards [post]
func (route Router) CreateWard(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.CreateWardRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.CreateWard(r.Context(), request)
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

// Ward godoc
// @Summary      Поиск подопечного
// @Description  Поиск подопечного по ID
// @Tags         Wards
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID подопечного"
// @Success      200  {object}  DatabaseServicev1.Ward
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/wards/{id} [get]
func (route Router) Ward(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindWardByIdRequest{Id: id}

	response, err := route.databaseService.FindWardById(r.Context(), request)
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

// DeleteWardByModel godoc
// @Summary      Удаление подопечного по модели
// @Description  Удаляет подопечного опираясь на всю сущность модели
// @Tags         Wards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ward body DatabaseServicev1.Ward false "Сущность подопечного"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/wards/deleteModel [post]
func (route Router) DeleteWardByModel(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.Ward)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.DeleteWardByModel(r.Context(), request)
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

// DeleteWardById godoc
// @Summary      Удаление подопечного
// @Description  Удаление подопечного по ID
// @Tags         Wards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID подопечного"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/wards/{id} [delete]
func (route Router) DeleteWardById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.DeleteWardByIdRequest{Id: id}

	response, err := route.databaseService.DeleteWardById(r.Context(), request)
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

// UpdateWard godoc
// @Summary      Обновление подопечного
// @Description  Обновление подопечного
// @Tags         Wards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ward body DatabaseServicev1.Ward false "Модель для обновления"
// @Success      200  {object}  DatabaseServicev1.Ward
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/wards [put]
func (route Router) UpdateWard(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.Ward)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.UpdateWard(r.Context(), request)
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

// FindWardDonations godoc
// @Summary      Извлечение пожертвований подопечного
// @Description  Извлечение пожертвований подопечного по его ID
// @Tags         Wards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Ward ID"
// @Success      200  {object}  DatabaseServicev1.DonationsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/wards/{id}/donations [get]
func (route Router) FindWardDonations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindWardDonationByIdRequest{Id: id}

	response, err := route.databaseService.FindWardDonationById(r.Context(), request)
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
