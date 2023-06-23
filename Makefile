EASYJSON_PATHS = ./internal/...

# DB
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

# Backend
.PHONY: build-image
build-image:
	DOCKER_BUILDKIT=1 docker build -t forum_backend .

.PHONY: run
run: build-image
	docker run --rm \
		-d \
        --memory 2G \
        --log-opt max-size=5M \
        --log-opt max-file=3 \
        --name f_backend \
		-p 5000:5000 \
	  	forum_backend

#		-d \

.PHONY: stop
stop:
	docker stop f_backend

.PHONY: logs
logs:
	docker logs -f f_backend


.PHONY: restart
restart: stop run

