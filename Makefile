.PHONY: run
run: down up

.PHONY: up
up:
	docker-compose -f docker-compose.yml --env-file ./.env up -d --build
	docker image prune --filter label=stage=builder --force

.PHONY: down
down:
	docker-compose -f docker-compose.yml --env-file ./.env down

.PHONY: test
test:
	docker-compose -f docker-compose.test.yml --env-file ./.env up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml --env-file ./.env down -v --rmi 'all'

.PHONY: clear
clear:
	docker-compose -f docker-compose.yml --env-file ./.env down -v
