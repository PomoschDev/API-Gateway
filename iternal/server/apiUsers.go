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
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  DatabaseServicev1.UsersResponse
// @Failure      400  {object}	HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users [get]
func (route Router) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := route.databaseService.Users(r.Context(), nil)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		code, errStr := utilities.GRPCErrToHttpErr(err)
		http.Error(w, errStr, code)
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
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}	HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id} [get]
func (route Router) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)
	user, err := route.databaseService.FindUserById(r.Context(), &DatabaseServicev1.FindUserByIdRequest{Id: id})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		code, errStr := utilities.GRPCErrToHttpErr(err)
		http.Error(w, errStr, code)
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
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body DatabaseServicev1.UpdateUserRequest true "Модель для обновления"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}	HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/{id} [put]
func (route Router) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := utilities.StrToUint(mux.Vars(r)["id"])

	updateUser := &DatabaseServicev1.UpdateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(updateUser); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updateUser.Id = id

	user, err := route.databaseService.UpdateUser(r.Context(), updateUser)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		code, errStr := utilities.GRPCErrToHttpErr(err)
		http.Error(w, errStr, code)
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
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body DatabaseServicev1.CreateUserRequest false "Сущность пользователя"
// @Success      200  {object}  DatabaseServicev1.CreateUserResponse
// @Failure      400  {object}	HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users [post]
func (route Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := new(DatabaseServicev1.CreateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	createdUser, err := route.databaseService.CreateUser(r.Context(), newUser)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		code, errStr := utilities.GRPCErrToHttpErr(err)
		http.Error(w, errStr, code)
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
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {string}  string    "ok"
// @Failure      400  {object}	HTTPError
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
		code, errStr := utilities.GRPCErrToHttpErr(err)
		http.Error(w, errStr, code)
		return
	}

	w.WriteHeader(int(response.Code))
}

// UserIsExists godoc
// @Summary      Проверка существует ли пользователь
// @Description  Проверка существует ли пользователь, проверка по номеру телефона
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        phone body DatabaseServicev1.UserIsExistsRequest true "Номер телефона"
// @Success      200  {object}  DatabaseServicev1.UserIsExistsResponse
// @Failure      400  {object}	HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/isexists [post]
func (route Router) UserIsExists(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.UserIsExistsRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	logger.Info("request: %+v", request)

	response, err := route.databaseService.UserIsExists(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		code, errStr := utilities.GRPCErrToHttpErr(err)
		http.Error(w, errStr, code)
		return
	}

	resp := struct {
		IsExists bool
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
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body DatabaseServicev1.IsRoleRequest true "Request"
// @Success      200  {object}  DatabaseServicev1.IsRoleResponse
// @Failure      400  {object}	HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/users/isrole [post]
func (route Router) UserIsRole(w http.ResponseWriter, r *http.Request) {
	request := new(DatabaseServicev1.IsRoleRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	logger.Info("request: %+v", request)

	response, err := route.databaseService.IsRole(r.Context(), request)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		code, errStr := utilities.GRPCErrToHttpErr(err)
		http.Error(w, errStr, code)
		return
	}

	str := utilities.ToJSON(response)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}
