EASYJSON_PATHS = ./internal/...

.PHONY: db-up
db-up:
	docker compose -f docker-compose.yml up -d db

.PHONY: db-stop
db-stop:
	docker compose -f docker-compose.yml stop db

.PHONY: db-down
db-down:
	docker compose -f docker-compose.yml stop db
	docker compose -f docker-compose.yml rm -f db

# easyjson
.PHONY: generate
generate:
	go generate ${EASYJSON_PATHS}
