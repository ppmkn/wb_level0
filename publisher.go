package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {
	// Настройки для подключения к NATS и NATS Streaming
	natsURL := "nats://localhost:4222"
	clusterID := "test-cluster"
	clientID := "my-unique-client2"
	subject := "sql-query" // Канал для SQL-запросов

	// сам SQL-запрос для отправки
	sqlQuery := `INSERT INTO orders (order_uid, track_number, entry)
                VALUES ('erdfg44dfgdfg', 'WBILMTESTTRACK', 'WBIL');`

	/* Стандартные данные что то вроде:
	   b111feb7b2b84b6test
	   WBILMTESTTRACK
	   WBIL */

	// Подключение к серверу NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS: %v", err)
	}
	defer nc.Close()

	// Подключение к NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Отправка SQL-запроса в канал
	err = sc.Publish(subject, []byte(sqlQuery))
	if err != nil {
		log.Printf("Ошибка при отправке SQL-запроса: %v", err)
	} else {
		log.Printf("SQL-запрос успешно отправлен:\n%s\n", sqlQuery)
	}
}
