package database

import (
	"fmt"
	"go-web-server-v2/models"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// // Формируем строку подключения
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	config.Config("DB_HOST"),
	// 	config.Config("DB_USER"),
	// 	config.Config("DB_PASSWORD"),
	// 	config.Config("DB_NAME"),
	// )

	// // Подключаемся к базе данных
	// var err error
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("Ошибка подключения к базе данных:", err)
	// }

	// Версия для RailWay
	// Берем DATABASE_URL из переменных окружения

	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL переменная окружения не установлена")
	}

	// Подключаемся к базе данных через GORM
	DB, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	if DB == nil {
		log.Fatal("db осталось nil", err)
	}
	defer DB.Close()

	// Проверка соединения с базой данных
	if err := DB.DB().Ping(); err != nil {
		log.Fatal("Не удалось установить соединение с базой данных:", err)
	} else {
		fmt.Println("Подключение к базе данных успешно установлено!")
	}

	// Автоматическая миграция (создаст таблицу users)
	DB.AutoMigrate(&models.User{})
}
