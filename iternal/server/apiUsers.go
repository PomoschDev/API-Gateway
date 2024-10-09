package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/utilities"
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

func (route Router) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := utilities.StrToUint(mux.Vars(r)["id"])

	updateUser := &DatabaseServicev1.UpdateUserRequest{}

	_ = updateUser

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
