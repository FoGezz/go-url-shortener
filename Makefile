build: 
	go build -o ./cmd/shortener ./cmd/shortener/ 
up:
	docker compose -f ./docker/docker-compose.yml up -d
down:
	docker compose -f ./docker/docker-compose.yml down