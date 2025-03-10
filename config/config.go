package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Загружаем конфиг из .env-файла
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic("Ошибка загрузки .env")
	}
}

func Config(key string) string {
	return os.Getenv(key)
}
