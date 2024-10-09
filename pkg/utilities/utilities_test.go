package utilities

import (
	"autoseller/backend/cmd/config"
	"autoseller/backend/pkg/logger"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"testing"
)

func initConfig() (*config.Config, error) {
	err := logger.New()
	if err != nil {
		return nil, err
	}
	pathToEnv := "C:\\Users\\Alex\\GolandProjects\\autoseller\\.env"

	err = godotenv.Load(pathToEnv)
	if err != nil {
		return nil, err
	}

	DeveloperID, _ := strconv.ParseInt(os.Getenv("developerId"), 10, 64)
	DatabasePort, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	RunPort, _ := strconv.Atoi(os.Getenv("RUN_PORT"))
	usrAnon := os.Getenv("USER_ANONYMOUS_MODE") == "1"  //1 true 0 false
	admAnon := os.Getenv("ADMIN_ANONYMOUS_MODE") == "1" //1 true 0 false
	alerts := os.Getenv("ALERTS") == "1"                //1 true 0 false
	channelID, _ := strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)
	comission, _ := strconv.ParseFloat(os.Getenv("COMISSION"), 32)

	Cfg := &config.Config{
		BotToken:               os.Getenv("BOT_TOKEN"),
		Developer:              os.Getenv("developer"),
		DeveloperID:            DeveloperID,
		DatabaseHost:           os.Getenv("DATABASE_HOST"),
		DatabasePort:           DatabasePort,
		DatabaseUser:           os.Getenv("DATABASE_USER"),
		DatabasePassword:       os.Getenv("DATABASE_PASSWORD"),
		Database:               os.Getenv("DATABASE"),
		RunPort:                RunPort,
		AESKey:                 os.Getenv("AESKEY"),
		RpcHost:                os.Getenv("RPC_HOST"),
		RpcUser:                os.Getenv("RPC_USER"),
		RpcPassword:            os.Getenv("RPC_PASSWORD"),
		WalletPassword:         os.Getenv("WALLET_PASSWORD"),
		UserAnonymous:          usrAnon,
		AdminAnonymous:         admAnon,
		BlockExplorerUrlFormat: os.Getenv("BLOCKEXPLORER_ULR"),
		JwtSecretKey:           []byte(os.Getenv("JWT_SECRET_KEY")),
		Alerts:                 alerts,
		ChannelID:              channelID,
		Comission:              comission,
		DeveloperWaller:        os.Getenv("DEVELOPER_WALLET"),
	}

	return Cfg, nil
}

func TestPercentage(t *testing.T) {
	type args struct {
		percent float64
		num     float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Тест1",
			args: args{
				percent: 5,
				num:     0.24,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Percentage(tt.args.percent, tt.args.num)
			t.Logf("Число: %f", got)
		})
	}
}

func TestGetCryptoKey(t *testing.T) {
	cfg, err := initConfig()
	if err != nil {
		t.Fatalf("Ошибка при инициализации конфига: %v", err)
		return
	}
	config.Cfg = cfg
	tests := []struct {
		name string
	}{
		{
			name: "Генерация ключа",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCryptoKey()
			t.Logf("Ключ: %s", got)
		})
	}
}

func TestGenerateQRCode(t *testing.T) {
	logger.New()
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Генерация qr-кода",
			args: args{str: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateQRCode(tt.args.str)
			if err != nil {
				t.Errorf("GenerateQRCode() error = %v", err)
				return
			}
			t.Logf("Путь: %s", got)
		})
	}
}
