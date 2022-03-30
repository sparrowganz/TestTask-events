# TestTask-events
Test Task for SayGames Company. Include web service for getting and saving events

---
## Task

Клиент посылает аналитические события пачками в виде POST запроса, в теле которого находятся сериализованные в JSON объекты, разделенные \n.
Каждый объект - это событие. Объект имеет такую структуру (в примере отформатированный JSON, в запросе - нет):

{
"client_time":"2020-12-01 23:59:00",
"device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
"device_os":"iOS 13.5.1",
"session":"ybuRi8mAUypxjbxQ",
"sequence":1,
"event":"app_start",
"param_int":0,
"param_str":"some text"
}

Клиенту достаточно вернуть ответ после чтения тела запроса.

Минимальный RPS на сервер - 200, среднее число событий в запросе - 30.

Задача - написать высокопроизводительный сервер для приема таких событий, их обогащения и вставки в БД.
Идеально будет развернуть эту систему на виртуальном сервачке.

Сервер пишем на Go, уместно используя его фичи - goroutines, channels.

В качестве БД используем ClickHouse, который надо поставить, создать в нем базу данных и табличку.

Код сервиса заливаем в любой git.

В процессе обогащения данных, добавляем в событие два поля - ip адрес клиента и время сервера:

    "ip":"8.8.8.8",
    "server_time":"2020-12-01 23:53:00"

Все, что не указано явно в задании, остается на твое усмотрение.

---

## Резюмируя
Данное решение не лишено огрехов и ошибок и носит демонстративный характер.

db/cache - что-то на подобие apache kafka (для сохранности евентов при завершении)

Текущий проект включает в себя:
- grace shutdown через errgroup
- Использование концепции ядра приложения
- DDD и clean arch
- ...

На моей локальной машине сервер выдержал предполагаемую нагрузку.

```Server Software:        
Server Hostname:        app
Server Port:            8080

Document Path:          /api/events
Document Length:        0 bytes

Concurrency Level:      500
Time taken for tests:   23.753 seconds
Complete requests:      20000
Failed requests:        0
Keep-Alive requests:    20000
Total transferred:      1980000 bytes
Total body sent:        132180000
HTML transferred:       0 bytes
Requests per second:    842.01 [#/sec] (mean)
Time per request:       593.815 [ms] (mean)
Time per request:       1.188 [ms] (mean, across all concurrent requests)
Transfer rate:          81.41 [Kbytes/sec] received
                        5434.44 kb/s sent
                        5515.84 kb/s total


Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2  14.4      0     139
Processing:     0  585 620.0    470    7065
Waiting:        0  584 619.9    469    7065
Total:          0  587 618.8    471    7065

Percentage of the requests served within a certain time (ms)
  50%    471
  66%    646
  75%    800
  80%    938
  90%   1384
  95%   1807
  98%   2417
  99%   2786
 100%   7065 (longest request)
```



