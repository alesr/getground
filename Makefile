.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yaml up -d --build

.PHONY: docker-down
docker-down: ## Stop docker containers and clear artefacts.
	docker-compose -f docker-compose.yaml down
	docker system prune 

.PHONY: db
db:
	@docker compose -f docker-compose.yaml up -d --build db

.PHONY: db_connect
db-connect:
	@mysql -h 127.0.0.1 -u root -p 

app:
	@docker compose -f docker-compose.yaml up --build app

test-unit:
	@go test -v -count=1 --race --cover --short ./...

test-integration: db
	@go test -v -count=1 --race --cover ./...
