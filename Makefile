build: 
	go build -o ./cmd/shortener ./cmd/shortener/ 
up:
	docker compose up -d
down:
	docker compose down