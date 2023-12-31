# Задание L0

При написании  **README** я предполагаю, что Вы открываете код через **IDE** и запускаете с помощью `go run`.

## Настройка проекта
1. Откройте проект через **IDE**.

2. Импортируйте базу данных `wb_database.sql` из папки `internal\database` в PostgreSQL.

3. Запустите **NATS Streaming Server**.

4. Запустить `main.go` с помощью команды `go run main.go`. В консоли появится сообщение: `Сервер запущен на порту 8080`.

## Использование HTTP-сервера

После запуска **HTTP-сервер** будет доступен по адресу `http://localhost:8080`. Вы можете отправить GET-запрос на этот адрес, указав в поле ввода ID заказа ( `OrderUID` ).

HTTP-сервер вернет найденные данные и запишет их в кэш.
Если ID не найден - сервер вернет страницу с ошибкой.

Все строки из базы данных обрезаются на сайте, если превышают лимит в **32** символа. Так не возникнет возможных проблем с отображением.


## Скрипт для публикации данных в канал

Необходимо запустить `sql_executor.go` с помощью `go run sql_executor.go`, а затем открыть сам скрипт публикации данных `publisher.go`.
Там можно найти следующую форму, для отправки новых данных:
```sql
    sqlQuery := `INSERT INTO orders (order_uid, track_number, entry)
                VALUES ('b111feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL');`
```
Редактируем значения и отправляем с помощью `go run publisher.go`
Появится сообщение, что данные отправлены.

Чтобы убедиться в этом наверняка, можно открыть `sql_executor.go` и найти в консоли следующее сообщение:
`SQL-запрос успешно выполнен: <данные запроса>`.

Если сообщение **отсутствует** - данные не дошли до базы данных.
