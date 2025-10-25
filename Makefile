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

make-topic:
	docker exec -it kafka kafka-topics.sh \
	--create \
	--topic orders \
	--bootstrap-server localhost:9092 \
	--partitions 1 \
	--replication-factor 1

	docker exec -it kafka kafka-topics.sh \
	--create \
	--topic orders_retry \
	--bootstrap-server localhost:9092 \
	--partitions 1 \
	--replication-factor 1

	docker exec -it kafka kafka-topics.sh \
	--create \
	--topic orders_dlq \
	--bootstrap-server localhost:9092 \
	--partitions 1 \
	--replication-factor 1
gen-mocks:
	@mockgen -source=internal/infrastructure/repo/repo.go -destination=mocks/mock_repo.go -package=mocks
	@mockgen -source=internal/infrastructure/cache/cache.go -destination=mocks/mock_cache.go -package=mocks