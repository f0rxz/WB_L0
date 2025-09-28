up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

ps:
	docker-compose ps

restart: down up

clean:
	docker-compose down -v
	docker volume prune -f

migrate-up:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/orderservice?sslmode=disable" up
