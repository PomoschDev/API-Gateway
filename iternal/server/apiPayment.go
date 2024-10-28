package server

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/logger"
	"apiGateway/pkg/token"
	"encoding/json"
	"net/http"
)

type PaymentRequest struct {
	ToWardId    uint64  `json:"toWardId"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// Payment godoc
// @Summary      Пожертвования
// @Description  Пожертвования
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payment body PaymentRequest true "Данные для оплаты"
// @Success      200  {object}  HTTPError
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/v1/payment [post]
func (route Router) Payment(w http.ResponseWriter, r *http.Request) {
	request := new(PaymentRequest)
	userId := r.Context().Value("user").(token.IUser).GetUserId()

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	if request.ToWardId <= 0 {
		SetHTTPError(w, "Поле \"ToWardId\" не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 {
		SetHTTPError(w, "Сумма пожертвования не может быть меньше или равно 0", http.StatusBadRequest)
		return
	}

	user, err := route.databaseService.FindUserById(r.Context(), &DatabaseServicev1.FindUserByIdRequest{Id: userId})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	cards, err := route.databaseService.FindUserCard(r.Context(), &DatabaseServicev1.FindUserCardRequest{Id: user.Id})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	ward, err := route.databaseService.FindWardById(r.Context(), &DatabaseServicev1.FindWardByIdRequest{Id: request.ToWardId})
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	if request.Description == "" || len(request.Description) == 0 {
		request.Description = ward.Want
	}

	donation := &DatabaseServicev1.CreateDonationsRequest{
		Title:  request.Description,
		Amount: float32(request.Amount),
		WardId: request.ToWardId,
		UserId: user.Id,
	}

	ward.Collected += donation.Amount

	_, err = route.databaseService.CreateDonations(r.Context(), donation)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	//TODO: потом изменить
	_ = cards

	_, err = route.databaseService.UpdateWard(r.Context(), ward)
	if err != nil {
		logger.Error("Ошибка при выполнении запроса: %v", err)
		SetGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
