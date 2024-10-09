package utilities

import (
	"autoseller/backend/cmd/config"
	"testing"
)

func TestEncryptAES(t *testing.T) {
	config.Init()
	type args struct {
		key     string
		message string
	}
	t.Logf("ключ: %s", config.Cfg.AESKey)
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Шифрование в AES",
			args: args{
				key:     config.Cfg.AESKey,
				message: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptAES(tt.args.key, tt.args.message)
			if err != nil {
				t.Errorf("EncryptAES() error = %v", err)
				return
			}
			t.Logf("Заширофанная строка: %s", got)
		})
	}
}

func TestDecryptAES(t *testing.T) {
	config.Init()
	type args struct {
		key     string
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Расшифровка",
			args: args{
				key:     config.Cfg.AESKey,
				message: "DzeeNUI/BB+8HBiYGdy1T6UZfcA=",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptAES(tt.args.key, tt.args.message)
			if err != nil {
				t.Errorf("DecryptAES() error = %v", err)
				return
			}
			t.Logf("Расшифрованная строка: %s", got)
		})
	}
}

func TestMD5(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Хэш логина",
			args: args{text: "111"},
		},
		{
			name: "Хэш пароля",
			args: args{text: "111"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MD5(tt.args.text)
			t.Logf("Результат: %s", got)
		})
	}
}
