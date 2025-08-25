# Makefile для WB Tech Demo Project
SHELL := /bin/bash

# Цвета для вывода
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: help run build-kafka-topic venv send-test clean

# Список команд: make без аргументов покажет справку
help:
	@echo ""
	@echo "${YELLOW}Доступные команды:${RESET}"
	@echo ""
	@echo "  ${GREEN}make run${RESET}           - Запустить всё: Docker, топик, venv, Go-сервис"
	@echo "  ${GREEN}make build${RESET}         - Только собрать Docker-инфраструктуру"
	@echo "  ${GREEN}make topic${RESET}         - Создать топик Kafka 'orders'"
	@echo "  ${GREEN}make venv${RESET}          - Создать и настроить виртуальное окружение Python"
	@echo "  ${GREEN}make send${RESET}          - Отправить тестовое сообщение в Kafka"
	@echo "  ${GREEN}make clean${RESET}         - Остановить всё и удалить данные"
	@echo ""

# Запуск всего проекта
run: build topic venv
	@echo "${GREEN}🚀 Запуск Go-сервиса...${RESET}"
	@echo "${YELLOW}Нажмите Ctrl+C для остановки${RESET}"
	@source venv/bin/activate && go run cmd/server/main.go

# Запуск Docker-контейнеров
build:
	@echo "${GREEN}🐳 Запуск Docker (PostgreSQL, Redis, Kafka)...${RESET}"
	docker-compose up -d

# Создание топика Kafka
topic:
	@echo "${GREEN}📌 Создание топика Kafka 'orders'...${RESET}"
	@docker exec -t kafka kafka-topics.sh --create \
		--topic orders \
		--bootstrap-server localhost:9092 \
		--partitions 1 \
		--replication-factor 1 2>/dev/null || \
		echo "${YELLOW}⚠️  Топик 'orders' уже существует или Kafka ещё не готов${RESET}"

# Настройка виртуального окружения Python
venv:
	@echo "${GREEN}🐍 Настройка Python виртуального окружения...${RESET}"
	@test -d venv || python3 -m venv venv
	@source venv/bin/activate && pip install --quiet kafka-python || \
	(echo "${GREEN}📦 Установка kafka-python...${RESET}" && \
	 source venv/bin/activate && pip install kafka-python)

# Отправка тестового сообщения
send:
	@echo "${GREEN}📤 Отправка тестового сообщения в Kafka...${RESET}"
	@source venv/bin/activate && python scripts/send_test_message.py

# Остановка и очистка
clean:
	@echo "${GREEN}🧹 Очистка: остановка Docker и удаление данных...${RESET}"
	docker-compose down -v
	@echo "${GREEN}🗑️  Удаление виртуального окружения (опционально)...${RESET}"
	# rm -rf venv
	@echo "${GREEN}✅ Очистка завершена${RESET}"