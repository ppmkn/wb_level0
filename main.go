package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// Определение структуры OrderData
type OrderData struct {
	OrderUID    string
	TrackNumber string
	Entry       string
}

// Настройки для NATS и PostgreSQL
const (
	natsURL       = "nats://localhost:4222"
	dbURL         = "postgres://postgres:root@localhost:5432/wb_database?sslmode=disable"
	natsClusterID = "test-cluster"
	natsClientID  = "my-unique-client"
	subject       = "sql-query" // Канал для SQL-запросов
	cacheLimit    = 100
)

var (
	natsConn *nats.Conn                    // подключение к NATS
	sc       stan.Conn                     // подключение к NATS Streaming
	db       *sql.DB                       // подключение к базе данных PostgreSQL
	cache    = make(map[string]*OrderData) // кэш
	mu       sync.RWMutex
)

func main() {
	fmt.Println("Сервер запущен на порту 8080") // Просто чтобы видеть, когда код запустился :)

	// Подключение к NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS: %v", err)
	}
	defer nc.Close()
	natsConn = nc

	// Подключение к NATS Streaming
	sc, err := stan.Connect(natsClusterID, natsClientID, stan.NatsConn(nc), stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
		log.Printf("Потеряно соединение с NATS Streaming: %v", err)
	}))
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Подключение к PostgreSQL
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	// Настройка HTTP сервера
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/result", resultHandler)
	http.ListenAndServe(":8080", nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request, message, errorMessage, orderID string) {
	// открыть error.html
	tmpl, err := template.ParseFiles("static/error.html")
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Создание данных для шаблона(сообщение об ошибке и введеный пользователем ID заказа)
	data := struct {
		Error   string
		OrderID string
	}{
		Error:   errorMessage,
		OrderID: limitString(orderID, 32),
	}

	// Установка HTTP-статуса "Not Found" (404)
	w.WriteHeader(http.StatusNotFound)

	// Выполнение шаблона с указанными данными и запись результатов в HTTP-ответ
	tmpl.Execute(w, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Если метод POST, значит была отправлена форма
		// Разбор данных формы из запроса
		r.ParseForm()
		orderID := r.FormValue("orderID")

		// Проверка кеша
		data, cached := getFromCache(orderID)
		if cached {
			http.Redirect(w, r, "/result?orderID="+orderID, http.StatusSeeOther)
			return
		}

		// Если нет в кеше, запросить из БД
		data, err := getFromDatabase(orderID)
		if err != nil {
			log.Printf("Ошибка при запросе данных из БД: %v", err)
			errorHandler(w, r, "Заказ не найден", "Заказ №", orderID)
			return
		}

		// Сохраняем данные в кеш
		saveToCache(orderID, data)

		// Перенаправляем на страницу результатов
		http.Redirect(w, r, "/result?orderID="+orderID, http.StatusSeeOther)
		return
	}
	// Если метод запроса не POST (например, GET), открываем страницу index.html
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	tmpl.Execute(w, nil)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра "orderID" из URL-запроса
	orderID := r.URL.Query().Get("orderID")

	// Проверка кеша
	data, cached := getFromCache(orderID)
	if !cached {
		// Если нет в кеше, запросить из БД
		data, err := getFromDatabase(orderID)
		if err != nil {
			http.Error(w, "Ошибка при запросе данных из БД", http.StatusInternalServerError)
			return
		}

		// Сохраняем данные в кеш
		saveToCache(orderID, data)
	}

	// Обрезаем данные до заданной длины ПЕРЕД отображением на странице
	trimData(data)
	// Открываем страницу result.html и передать данные для отображения
	tmpl := template.Must(template.ParseFiles("static/result.html"))
	tmpl.Execute(w, data)
}

func getFromDatabase(orderID string) (*OrderData, error) {
	// Выполняем SQL-запрос к БД
	row := db.QueryRow("SELECT order_uid, track_number, entry FROM orders WHERE order_uid = $1", orderID)

	// Создаем структуру для хранения данных о заказе
	var data OrderData

	// Извлекаем данные из результата запроса и сохраняем их в структуре
	err := row.Scan(&data.OrderUID, &data.TrackNumber, &data.Entry)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Заказ не найден")
		}
		// если произошла другая ошибка
		return nil, err
	}

	return &data, nil
}

func getFromCache(orderID string) (*OrderData, bool) {
	// Получаем мьютекс для чтения из кеша
	mu.RLock()
	defer mu.RUnlock() // Убеждаемся, что мьютекс освобождается после завершения функции

	// Попытка извлечь данные из кеша для указанного orderID
	data, cached := cache[orderID]

	// Возвращаем данные (если они есть) и флаг, указывающий, были ли данные найдены в кеше
	return data, cached
}

func saveToCache(orderID string, data *OrderData) {
	// Получаем мьютекс для записи в кеш
	mu.Lock()
	defer mu.Unlock() // Убеждаемся, что мьютекс освобождается после завершения функции

	if len(cache) >= cacheLimit {
		// Если кеш достиг максимального размера, удаляем 1 элемент
		for key := range cache {
			delete(cache, key)
			break
		}
	}

	// Сохраняем данные в кеш для указанного orderID
	cache[orderID] = data
}

// Функции для обрезания слишком длинных строк, чтобы в случаи чего на фронте всё выглядело адекватно
func limitString(input string, maxLength int) string { // Обрезание данных до 32 символов
	if len(input) > maxLength {
		return input[:maxLength] + "..." // 3 точки в конце это знак для пользователя, что текст обрезан
	}
	return input
}

func trimData(data *OrderData) { // Для вывода в index.html
	data.OrderUID = limitString(data.OrderUID, 32)
	data.TrackNumber = limitString(data.TrackNumber, 32)
	data.Entry = limitString(data.Entry, 32)
}
