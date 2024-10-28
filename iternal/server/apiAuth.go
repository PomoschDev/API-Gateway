package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/token"
	"apiGateway/pkg/utilities"
	"encoding/json"
	"github.com/nyaruka/phonenumbers"
	"net/http"
	"regexp"
	"unicode"
)

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegistrationRequest struct {
	Email    string                                  `json:"email,omitempty"`
	Username string                                  `json:"username,omitempty"`
	Password string                                  `json:"password,omitempty"`
	Phone    string                                  `json:"phone,omitempty"`
	Card     *DatabaseServicev1.CreateCardRequest    `json:"card,omitempty"`
	Company  *DatabaseServicev1.CreateCompanyRequest `json:"company,omitempty"`
	Type     uint64                                  `json:"type"`
}

// Login godoc
// @Summary      Авторизация
// @Description  Авторизация пользователя
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user body LoginRequest true "Данные для авторизации пользователя"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  HTTPError
// @Failure      401  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/auth/login [post]
func (route Router) Login(w http.ResponseWriter, r *http.Request) {
	request := new(LoginRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	if request.Phone == "" || len(request.Phone) == 0 {
		logger.Error("Поле phone пустое")
		SetHTTPError(w, "Поле \"Phone\" не может быть пустым", http.StatusBadRequest)
		return
	}

	if request.Password == "" || len(request.Password) == 0 {
		logger.Error("Поле phone пустое")
		SetHTTPError(w, "Поле \"Password\" не может быть пустым", http.StatusBadRequest)
		return
	}

	user, err := route.databaseService.FindUserByPhone(r.Context(),
		&DatabaseServicev1.FindUserByPhoneRequest{Phone: request.Phone})
	if err != nil {
		SetGRPCError(w, err)
		return
	}

	comparePasswordRequest := new(DatabaseServicev1.ComparePasswordRequest)
	err = utilities.Transformation(request, comparePasswordRequest)
	if err != nil {
		SetHTTPError(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}

	responseComparePassword, err := route.databaseService.ComparePassword(r.Context(), comparePasswordRequest)
	if err != nil {
		SetGRPCError(w, err)
		return
	}

	if !responseComparePassword.Accessory {
		SetHTTPError(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	jwtToken, err := token.CreateToken(user, route.cfg)
	if err != nil {
		logger.Error("Ошибка при создании JWT токена: %v", err)
		SetHTTPError(w, "Ошибка при создании JWT токена", http.StatusInternalServerError)
		return
	}

	response := &LoginResponse{Token: jwtToken}

	str := utilities.ToJSON(response)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// Registration godoc
// @Summary      Регистрация пользователя
// @Description  Регистрация нового пользователя пользователя
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user body RegistrationRequest false "Данные для регистрации"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/auth/registration [post]
func (route Router) Registration(w http.ResponseWriter, r *http.Request) {
	registrationRequest := new(RegistrationRequest)

	if err := json.NewDecoder(r.Body).Decode(registrationRequest); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	logger.Info("registrationRequest: %s", utilities.ToJSON(registrationRequest))

	phone, err := phonenumbers.Parse(registrationRequest.Phone, "RU")
	if err != nil {
		SetHTTPError(w, "Ошибка при извлечении номера телефона", http.StatusBadRequest)
		return
	}

	if !phonenumbers.IsValidNumber(phone) {
		SetHTTPError(w, "Неверный формат номера телефона", http.StatusBadRequest)
		return
	}

	// Регулярное выражение для проверки формата email
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Проверка email
	if !re.MatchString(registrationRequest.Email) {
		SetHTTPError(w, "Неверный формат почты", http.StatusBadRequest)
		return
	}

	if registrationRequest.Password == "" || len(registrationRequest.Password) == 0 {
		SetHTTPError(w, "Пароль не может быть пустым", http.StatusBadRequest)
		return
	}

	if len(registrationRequest.Password) < 8 {
		SetHTTPError(w, "Пароль не может быть меньше 8 символов", http.StatusBadRequest)
		return
	}

	if !hasUppercase(registrationRequest.Password) {
		SetHTTPError(w, "Должна быть хотя бы одна заглавная буква", http.StatusBadRequest)
		return
	}

	newUser := &DatabaseServicev1.CreateUserRequest{
		Email:     registrationRequest.Email,
		Username:  registrationRequest.Username,
		Password:  registrationRequest.Password,
		Phone:     registrationRequest.Phone,
		Card:      []*DatabaseServicev1.CreateCardRequest{registrationRequest.Card},
		Role:      "user",
		Company:   registrationRequest.Company,
		Type:      registrationRequest.Type,
		Donations: nil,
	}

	respService, err := route.databaseService.CreateUser(r.Context(), newUser)
	if err != nil {
		SetGRPCError(w, err)
		return
	}

	jwtToken, err := token.CreateToken(respService, route.cfg)
	if err != nil {
		logger.Error("Ошибка при создании JWT токена: %v", err)
		SetHTTPError(w, "Ошибка при создании JWT токена", http.StatusInternalServerError)
		return
	}

	response := &LoginResponse{Token: jwtToken}

	str := utilities.ToJSON(response)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("Ошибка при формировании ответа: %v", err)
	}
}
