package utilities

import (
	"apiGateway/pkg/logger"
	"bytes"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"net/http"
	"strconv"
)

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "   ")
	if err != nil {
		return in
	}
	return out.String()
}

// ToJSON - конвертирует объект в JSON строку
func ToJSON(object any) string {
	jsonByte, err := json.Marshal(object)
	if err != nil {
		logger.Error("Ошибка при получении JSON: %v", err)
	}
	n := len(jsonByte)             //Find the length of the byte array
	result := string(jsonByte[:n]) //convert to string

	return jsonPrettyPrint(result)
}

// StrToUint - Конвертирует строку в uint
func StrToUint(s string) uint64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Error("%v", err)
		return 0
	}
	return uint64(i)
}

// RandInt - возвращает случайное число от min до max
func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// GenerateRandomString - генерирует случайный набор символов (англ алфавит, case uppercase + символ _ и цифры от 0 до 9)
func GenerateRandomString(length int) string {
	alphabet := "QOS4rT08Dm7dZVOPwucfM2haFiNyEjBK3UtC9IqY_lzv6XpWgWsAJebG5H1RxnLbK"

	var result = make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(result)
}

// Transformation - преобразование одной модели в другую
func Transformation(forModel any, toModel any) error {
	encodedJsonModelBytes, err := json.Marshal(forModel)
	if err != nil {
		return err
	}

	err = json.Unmarshal(encodedJsonModelBytes, toModel)
	if err != nil {
		return err
	}

	return nil
}

func GRPCErrToHttpErr(grpcErr error) (int, string) {
	statusCode, _ := status.FromError(grpcErr)
	logger.Info("Status code: %+v", statusCode)
	errStr := statusCode.Message()
	httpCode := 0

	switch statusCode.Code() {
	case codes.OK:
		httpCode = http.StatusOK
	case codes.Canceled:
		httpCode = http.StatusConflict
	case codes.Unknown:
		httpCode = http.StatusInternalServerError
	case codes.InvalidArgument:
		httpCode = http.StatusBadRequest
	case codes.DeadlineExceeded:
		httpCode = http.StatusGatewayTimeout
	case codes.NotFound:
		httpCode = http.StatusNotFound
	case codes.AlreadyExists:
		httpCode = http.StatusConflict
	case codes.PermissionDenied:
		httpCode = http.StatusForbidden
	case codes.ResourceExhausted:
		httpCode = http.StatusTooManyRequests
	case codes.FailedPrecondition:
		httpCode = http.StatusBadRequest
	case codes.Aborted:
		httpCode = http.StatusConflict
	case codes.OutOfRange:
		httpCode = http.StatusBadRequest
	case codes.Unimplemented:
		httpCode = http.StatusNotImplemented
	case codes.Internal:
		httpCode = http.StatusInternalServerError
	case codes.Unavailable:
		httpCode = http.StatusServiceUnavailable
	case codes.DataLoss:
		httpCode = http.StatusInternalServerError
	case codes.Unauthenticated:
		httpCode = http.StatusUnauthorized
	default:
		httpCode = http.StatusInternalServerError
	}

	return httpCode, errStr
}
