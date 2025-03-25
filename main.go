package main

import (
	"fmt"
	"go-web-server-v2/config"
	"go-web-server-v2/database"
	"go-web-server-v2/handlers"
	"go-web-server-v2/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Загружаем конфиг
	config.LoadConfig()

	// Подключаем базу данных
	database.ConnectDB()

	// var users []models.User
	// result := database.DB.Find(&users)
	// if result.Error != nil {
	// 	log.Fatal("Ошибка запроса:", result.Error)
	// }

	// fmt.Println("Список пользователей:")
	// for _, user := range users {
	// 	fmt.Printf("ID: %d, Username: %s Password: %s\n", user.ID, user.Username, user.Password)
	// }

	// // Удаление всех записей
	// database.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE;")

	// fmt.Println("✅ База данных очищена!")

	// Создаём маршруты
	r := mux.NewRouter()

	// Блокиратор
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	imgHandler := http.StripPrefix("/img", http.FileServer(http.Dir("./img")))

	// Какой-то FileServer
	r.PathPrefix("/static/").Handler(middleware.BlockFileDownload(staticHandler))
	r.PathPrefix("/img/").Handler((imgHandler))

	// GET-запросы
	r.HandleFunc("/", handlers.WithTemplateFile("templates/index.html", handlers.HomeHandler)).Methods("GET")
	r.HandleFunc("/register", handlers.WithTemplateFile("templates/register.html", handlers.HomeHandler)).Methods("GET")
	r.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")

	// POST-запросы
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Защищённые роуты (требуют JWT)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)
	api.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// Запускаем сервер
	port := config.Config("PORT")
	fmt.Println("Сервер запущен на порту:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
