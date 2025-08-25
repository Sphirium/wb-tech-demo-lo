### WB Tech Demo Project ###

Простой Go-сервис для приёма, хранения и отображения заказов из Kafka.  
Использует PostgreSQL, Redis, Kafka и веб-интерфейс.

## Ссылка на Google-диск с видеодемонстрацией работы пет-проекта:

https://drive.google.com/file/d/1cCqTwHFkZh7sLZmuSXl7mNWMZn039W7n/view?usp=sharing

## Основной стек ###

- Go 1.24.6
- Kafka
- PostgreSQL — основное хранилище
- Redis — кеширование
- Веб-интерфейс (статический HTML)

---

## Быстрый старт

1. **Запусти инфраструктуру и сервис:**

* `make run`


2. **Отправь тестовое сообщение:**

* `make send` (в другом терминале)


3. **Открой браузер:**

http://localhost:8081


4. Введи order_uid из лога — получи JSON заказа


## Основные команды:

* `make run` - запустить всё и запустить Go-сервис
* `make send` - отправить тестовый заказ в Kafka
* `make build` - запустить только Docker (без Go)
* `make clean` - остановить всё и удалить данные

## Запуск интеграционных тестов:

* `go test -v ./test/integration/`


## Требования:

- Docker
- Go 1.22+
- Python 3.8+