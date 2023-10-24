package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const (
	dbURL         = "postgres://postgres:root@localhost:5432/wb_database?sslmode=disable"
	natsURL       = "nats://localhost:4222"
	natsClusterID = "test-cluster"
	subject       = "sql-query"
)

var (
	db *sql.DB
)

func main() {
	fmt.Println("Скрипт запущен")

	// Подключение к NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS: %v", err)
	}
	defer nc.Close()

	// Подключение к NATS Streaming
	sc, err := stan.Connect(natsClusterID, "sql-processor", stan.NatsConn(nc))
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

	// Обработчик для сообщений из NATS
	messageHandler := func(msg *stan.Msg) {
		// Получаем SQL-запрос из сообщения
		sqlQuery := string(msg.Data)

		// Выполняем SQL-запрос в базе данных
		_, err := db.Exec(sqlQuery)
		if err != nil {
			log.Printf("Ошибка выполнения SQL-запроса: %v", err)
		} else {
			log.Printf("SQL-запрос успешно выполнен:\n%s\n", sqlQuery)
		}
	}

	// Подписываемся на канал NATS Streaming
	sub, err := sc.Subscribe(subject, messageHandler)
	if err != nil {
		log.Fatalf("Ошибка подписки на канал: %v", err)
	}
	defer sub.Unsubscribe()

	// Ожидаем завершения приложения при получении сигнала завершения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
}
