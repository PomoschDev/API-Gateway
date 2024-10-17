package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/utilities"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// GetUsers godoc
// @Summary      Список всех пользователей
// @Description  Массив пользователей в базе данных
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  DatabaseServicev1.UsersResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users [get]
func (route Router) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := route.databaseService.Users(r.Context(), nil)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(users)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// GetUser godoc
// @Summary      Поиск пользователя
// @Description  Поиск пользователя по ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id} [get]
func (route Router) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)
	user, err := route.databaseService.FindUserById(r.Context(), &DatabaseServicev1.FindUserByIdRequest{Id: id})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(user)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// UpdateUser godoc
// @Summary      Обновление пользователя
// @Description  Обновление сущности пользователя
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body DatabaseServicev1.UpdateUserRequest true "Модель для обновления"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id} [put]
func (route Router) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := utilities.StrToUint(mux.Vars(r)["id"])

	updateUser := &DatabaseServicev1.UpdateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(updateUser); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	updateUser.Id = id

	user, err := route.databaseService.UpdateUser(r.Context(), updateUser)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(user)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// CreateUser godoc
// @Summary      Создание пользователя
// @Description  Создание новой сущности пользователя
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body DatabaseServicev1.CreateUserRequest false "Сущность пользователя"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users [post]
func (route Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := new(DatabaseServicev1.CreateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	createdUser, err := route.databaseService.CreateUser(r.Context(), newUser)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(createdUser)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// DeleteUserByID godoc
// @Summary      Удаление пользователя
// @Description  Удаление пользователя по ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id} [delete]
func (route Router) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	request := &DatabaseServicev1.DeleteUserByIdRequest{Id: id}

	response, err := route.databaseService.DeleteUserById(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	str := utilities.ToJSON(response)
	_, _ = w.Write([]byte(str))
}

// UserIsExists godoc
// @Summary      Проверка существует ли пользователь
// @Description  Проверка существует ли пользователь, проверка по номеру телефона
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        phone body DatabaseServicev1.UserIsExistsRequest true "Номер телефона"
// @Success      200  {object}  DatabaseServicev1.UserIsExistsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/isExists [post]
func (route Router) UserIsExists(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.UserIsExistsRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	logger.Info("request: %+v", request)

	response, err := route.databaseService.UserIsExists(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	resp := struct {
		IsExists bool `json:"isExists"`
	}{
		IsExists: response.IsExists,
	}

	str := utilities.ToJSON(resp)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// UserIsRole godoc
// @Summary      Проверяет принадлежность к роли
// @Description  Проверяет пользователя на принадлежность к определенной роли
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body DatabaseServicev1.IsRoleRequest true "Request"
// @Success      200  {object}  DatabaseServicev1.IsRoleResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/isRole [post]
func (route Router) UserIsRole(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.IsRoleRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	logger.Info("request: %+v", request)

	response, err := route.databaseService.IsRole(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	resp := struct {
		Accessory bool `json:"accessory"`
	}{
		Accessory: response.Accessory,
	}

	str := utilities.ToJSON(resp)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// FindUserByEmail godoc
// @Summary      Поиск пользователя по email
// @Description  Поиск пользователя по его email
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        email query string true "Email" Format(email)
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/ [get]
func (route Router) FindUserByEmail(w http.ResponseWriter, r *http.Request) {
	request := &DatabaseServicev1.FindUserByEmailRequest{Email: r.URL.Query().Get("email")}

	if request.Email == "" || len(request.Email) == 0 {
		logger.Error("Email не может быть пустым")
		SetHTTPError(w, "Поле \"email\" не может быть пустым", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.FindUserByEmail(r.Context(), request)
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

// ComparePassword godoc
// @Summary      Сравнение вводимого пароля от пользователя
// @Description  Сравнивает пароль что ввел пользователь, с тем что есть в базе данных у его аккаунта
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body DatabaseServicev1.ComparePasswordRequest true "Данные пользователя"
// @Success      200  {object}  DatabaseServicev1.ComparePasswordResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/comparePassword [post]
func (route Router) ComparePassword(w http.ResponseWriter, r *http.Request) {
	request := &DatabaseServicev1.ComparePasswordRequest{}

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.ComparePassword(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	resp := struct {
		Accessory bool `json:"accessory"`
	}{
		Accessory: response.Accessory,
	}

	str := utilities.ToJSON(resp)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// ChangeUserType godoc
// @Summary      Меняет тип пользователя
// @Description  Обновление типа пользователя (0 - юридическое лицо, 1 - физическое лицо)
// @Tags         Users
// @Accept       mpfd
// @Produce      json
// @Param        id   path      uint64  true  "ID пользователя"
// @Param        type formData  uint64  true  "Type поле пользователя, 0 - юридическое лицо, 1 - физическое лицо"
// @Success      200  {object}  DatabaseServicev1.ChangeUserTypeResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id} [patch]
func (route Router) ChangeUserType(w http.ResponseWriter, r *http.Request) {
	id := utilities.StrToUint(mux.Vars(r)["id"])

	if err := r.ParseMultipartForm(1); err != nil {
		logger.Error("Ошибка при парсе формы: %v", err)
		SetHTTPError(w, "Слишком длинное тело запроса", http.StatusBadRequest)
		return
	}

	userType := utilities.StrToUint(r.FormValue("type"))

	if userType < 0 {
		SetHTTPError(w, "Поле \"type\" не может быть меньше 0", http.StatusBadRequest)
		return
	}

	if userType > 1 {
		SetHTTPError(w, "Поле \"type\" не может быть больше 1", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.ChangeUserTypeRequest{
		Id:   id,
		Type: userType,
	}

	response, err := route.databaseService.ChangeUserType(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	resp := struct {
		Accessory bool `json:"accessory"`
	}{
		Accessory: response.Accessory,
	}

	str := utilities.ToJSON(resp)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// FindUserByPhone godoc
// @Summary      Поиск пользователя по номеру телефона
// @Description  Поиск пользователя по его phone
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        phone query string true "Phone"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/ [get]
func (route Router) FindUserByPhone(w http.ResponseWriter, r *http.Request) {
	request := &DatabaseServicev1.FindUserByPhoneRequest{Phone: r.URL.Query().Get("phone")}

	if request.Phone == "" || len(request.Phone) == 0 {
		logger.Error("Phone не может быть пустым")
		SetHTTPError(w, "Поле \"phone\" не может быть пустым", http.StatusBadRequest)
		return
	}

	logger.Info("request: %+v", request)

	response, err := route.databaseService.FindUserByPhone(r.Context(), request)
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

// FindUserCompany godoc
// @Summary      Извлечение компании пользователя
// @Description  Извлечение компании пользователя по его ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  DatabaseServicev1.Company
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id}/company [get]
func (route Router) FindUserCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindUserCompanyRequest{Id: id}

	response, err := route.databaseService.FindUserCompany(r.Context(), request)
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

// FindUserDonations godoc
// @Summary      Извлечение пожертвований пользователя
// @Description  Извлечение пожертвований пользователя по его ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  DatabaseServicev1.FindUserDonationsResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id}/donation [get]
func (route Router) FindUserDonations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindUserDonationsRequest{Id: id}

	response, err := route.databaseService.FindUserDonations(r.Context(), request)
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

// FindUserCard godoc
// @Summary      Извлечение карт пользователя
// @Description  Извлечение карт пользователя по его ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  DatabaseServicev1.FindUserCardResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id}/card [get]
func (route Router) FindUserCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Поле \"ID\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	request := &DatabaseServicev1.FindUserCardRequest{Id: id}

	response, err := route.databaseService.FindUserCard(r.Context(), request)
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

// AddCardToUser godoc
// @Summary      Добавляет банковскую карту пользователю
// @Description  Добавляет банковскую карту пользователю, поле userId это ID пользователя в базе данных, которому будем добавлять карту
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        card body DatabaseServicev1.AddCardToUserRequest false "Сущность банковской карты"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/addCard [post]
func (route Router) AddCardToUser(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.AddCardToUserRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.AddCardToUser(r.Context(), request)
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

// DeleteUserByModel godoc
// @Summary      Удаление пользователя по модели
// @Description  Удаляет пользователя опираясь на всю сущность модели
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body DatabaseServicev1.DeleteUserByModelRequest false "Модель пользователя"
// @Success      200  {object}  DatabaseServicev1.HTTPCodes
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/deleteModel [post]
func (route Router) DeleteUserByModel(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.DeleteUserByModelRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	response, err := route.databaseService.DeleteUserByModel(r.Context(), request)
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
