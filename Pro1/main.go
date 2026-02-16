package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var config Config
var users []User

func main() {
	// Загружаем конфиг
	loadConfig()

	// Инициализируем данные
	users = []User{
		{ID: 1, Name: "Анна", Age: 25},
		{ID: 2, Name: "Петр", Age: 30},
	}

	// Настройка маршрутов
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/", handleUserByID)

	port := ":" + config.Port
	fmt.Printf("Сервер '%s' запущен на порту %s\n", config.Name, port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func loadConfig() {
	// Проверяем наличие файла конфигурации
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		// Создаем конфиг по умолчанию
		config = Config{Port: "8080", Name: "Pro1"}
		saveConfig()
		return
	}

	// Читаем конфиг
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Ошибка чтения конфига:", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Ошибка парсинга конфига:", err)
	}
}

func saveConfig() {
	data, _ := json.MarshalIndent(config, "", " ")
	ioutil.WriteFile("config.json", data, 0644)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Добро пожаловать в %s!", config.Name)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(users)
	} else if r.Method == "POST" {
		var newUser User
		json.NewDecoder(r.Body).Decode(&newUser)
		newUser.ID = len(users) + 1
		users = append(users, newUser)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}

func handleUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем ID из URL (например, /users/1)
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		http.NotFound(w, r)
		return
	}

	// Преобразуем ID в число
	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Ищем пользователя по ID
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	// Если пользователь не найден
	http.Error(w, "Пользователь не найден", http.StatusNotFound)
}
