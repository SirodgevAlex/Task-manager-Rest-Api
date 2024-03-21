Привет!

1) Запускаем db.sql

   ```pgsql
   psql db.sql
   ```

   у себя я делал не так, ибо у меня там pgAdmin

   ```pgsql
   psql -U postgres -d db -f db.sql -W 
   ```

   postgres - название пользователя в pgAdmin, у меня потом пароль попросило от пользователя
2) Запускаем main.go

   ```go
   go run main.go
   ```
3) Все, мы все запустили, можно делать сами запросы
4) Запрос для создания пользователя

   ```bash
   curl -i -X POST \
     http://localhost:8080/users \
     -H 'Content-Type: application/json' \
     -d '{
       "Name": "John Glick",
       "Balance": 100
   }'
   ```

   получим такой результат

   ```json
   HTTP/1.1 201 Created
   Date: Mon, 19 Mar 2024 00:46:52 GMT
   Content-Length: 41
   Content-Type: text/plain; charset=utf-8

   {"Id":3,"Name":"John Glick","Balance":100}
   ```
5) Запрос для создания quest

```bash
   curl -i -X POST \
     http://localhost:8080/quests \
     -H 'Content-Type: application/json' \
     -d '{
       "Name": "quest1",
       "Cost": 20
   }'
```

   получим такой результат

```json
HTTP/1.1 201 Created
Date: Wed, 20 Mar 2024 06:09:03 GMT
Content-Length: 34
Content-Type: text/plain; charset=utf-8

{"Id":1,"Name":"quest1","Cost":20}
```

6. Запрос для уведомления сервиса о выполенном задании

   ```bash
   curl -i -X POST \
     http://localhost:8080/complete \
     -H 'Content-Type: application/json' \
     -d '{
       "userId": 1,
       "questId": 4
   }'
   ```

   получим такой результат

   ```bash
   HTTP/1.1 200 OK
   Date: Mon, 19 Mar 2024 09:59:26 GMT
   Content-Length: 32
   Content-Type: text/plain; charset=utf-8

   {"Id":0,"UserId":1,"QuestId":4}
   ```

   предусмотрена обработка краевых случаев
7. Запрос для получения истории пользователя

   ```bash
   curl -i -X GET http://localhost:8080/history/1
   ```

   единичка на конце это userId

P.S. я долго пытался сделать через докер - не смог, ибо так и не понял, где прописывать запуск db.sql, чтоб, когда начнем совершать события, данные клались в соответсвующую бд, а бд не были созданы, потому что не не смог запустить db.sql.
