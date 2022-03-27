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
Таких как: 
- Черезмерное использование приведения типов.
- Одинарные запросы в clickhouse (max 100 дальше ошибках)
- итд

В  начале было в планах сделать все через кластер kafka и балансировку на воркеры. 
Но таска же на golang разработчика.😂

В этом тестовом задании я постарался показать свои умения организовывать и писать код, а так же возможность мыслить нестандартно.

На моей локальной машине сервер выдержал предполагаемую нагрузку. (Не считая потери нескольих евентов)

```Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /api/events
Document Length:        0 bytes

Concurrency Level:      300
Time taken for tests:   13.432 seconds
Complete requests:      5000
Failed requests:        10
   (Connect: 10, Receive: 0, Length: 0, Exceptions: 0)
Total transferred:      375000 bytes
Total body sent:        32985000
HTML transferred:       0 bytes
Requests per second:    372.24 [#/sec] (mean)
Time per request:       805.929 [ms] (mean)
Time per request:       2.686 [ms] (mean, across all concurrent requests)
Transfer rate:          27.26 [Kbytes/sec] received
                        2398.12 kb/s sent
                        2425.39 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2  12.2      0     102
Processing:     8  785 350.3    750    2338
Waiting:        7  735 322.2    702    2314
Total:          9  788 348.6    752    2338

Percentage of the requests served within a certain time (ms)
  50%    752
  66%    930
  75%   1002
  80%   1070
  90%   1228
  95%   1385
  98%   1643
  99%   1735
 100%   2338 (longest request)
```
Полная обработка данных заняла примерно 1мин (долетали евенты в базу)
После теста появились 149939 новых евента в базе



