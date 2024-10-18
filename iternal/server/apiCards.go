package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/utilities"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Cards godoc
// @Summary      Список всех банковских карт
// @Description  Массив банковских карт в базе данных
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Success      200  {object}  DatabaseServicev1.CardsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/cards [get]
func (route Router) Cards(w http.ResponseWriter, r *http.Request) {
	response, err := route.databaseService.Cards(r.Context(), nil)
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

// Card godoc
// @Summary      Поиск банковской карты
// @Description  Поиск банковской карты по ID
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Card ID"
// @Success      200  {object}  DatabaseServicev1.Card
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/cards/{id} [get]
func (route Router) Card(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindCardByIdRequest{Id: id}

	response, err := route.databaseService.FindCardById(r.Context(), request)
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

// CreateCard godoc
// @Summary      Создание банковской карты пользователя
// @Description  Создание банковской карты пользователя
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        card body DatabaseServicev1.CreateCardRequest false "Сущность банковской карты"
// @Success      200  {object}  DatabaseServicev1.Card
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/cards [post]
func (route Router) CreateCard(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.CreateCardRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	if request.GetUserId() <= 0 {
		SetHTTPError(w, "Поле \"UserID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	_, err := route.databaseService.FindUserById(r.Context(), &DatabaseServicev1.FindUserByIdRequest{Id: request.
		GetUserId()})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	response, err := route.databaseService.CreateCard(r.Context(), request)
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

// DeleteCardByModel godoc
// @Summary      Удаление банковской карты по модели
// @Description  Удаляет банковскую карту опираясь на всю сущность модели
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        card body DatabaseServicev1.Card false "Модель банковской карты"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/cards/deleteModel [post]
func (route Router) DeleteCardByModel(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.Card)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.DeleteCardByModel(r.Context(), request)
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

// DeleteCardById godoc
// @Summary      Удаление банковской карты пользователя
// @Description  Удаление банковской карты пользователя по ID
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID банковской карты"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/cards/{id} [delete]
func (route Router) DeleteCardById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	request := &DatabaseServicev1.DeleteCardByIdRequest{Id: id}

	response, err := route.databaseService.DeleteCardById(r.Context(), request)
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

// UpdateCard godoc
// @Summary      Обновление банковской карты
// @Description  Обновление банковской карты
// @Tags         Cards
// @Accept       json
// @Produce      json
// @Param        id path int true "ID банковской карты"
// @Param        card body DatabaseServicev1.UpdateUserCardRequest1 true "Модель для обновления"
// @Success      200  {object}  DatabaseServicev1.UpdateUserCardResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/cards/{id} [put]
func (route Router) UpdateCard(w http.ResponseWriter, r *http.Request) {
	id := utilities.StrToUint(mux.Vars(r)["id"])

	request := new(DatabaseServicev1.UpdateUserCardRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	request.Id = id

	response, err := route.databaseService.UpdateCard(r.Context(), request)
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
