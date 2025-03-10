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

	// Создаём маршруты
	r := mux.NewRouter()

	// Блокиратор
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	imgHandler := http.StripPrefix("/img", http.FileServer(http.Dir("./img")))

	// Какой-то FileServer
	r.PathPrefix("/static/").Handler(middleware.BlockFileDownload(staticHandler))
	r.PathPrefix("/img/").Handler((imgHandler))

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
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
