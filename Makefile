# Makefile
COMPOSE := docker compose

.PHONY: infra-up infra-down infra-restart infra-logs infra-ps infra-build infra-clean

infra-up:
	$(COMPOSE) up -d postgis valhalla

infra-down:
	$(COMPOSE) down

infra-restart:
	$(COMPOSE) restart postgis valhalla

infra-logs:
	$(COMPOSE) logs -f postgis valhalla

infra-ps:
	$(COMPOSE) ps

infra-build:
	$(COMPOSE) up --build -d

infra-clean:
	$(COMPOSE) down -v --remove-orphans
	
.PHONY: up down build logs ps clean

up:
	docker-compose up -d

build:
	docker-compose up --build

down:
	docker-compose down

logs:
	docker-compose logs -f

ps:
	docker-compose ps

clean:
	docker-compose down -v
