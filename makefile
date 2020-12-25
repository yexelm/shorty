run: down up

up:
	go clean -cache
	docker-compose -f docker-compose.yml up -d --build

down:
	docker-compose -f docker-compose.yml down

test:
	go clean -cache
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down

clear:
	docker-compose -f docker-compose.yml down -v
