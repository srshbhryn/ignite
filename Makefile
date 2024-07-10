test:
	docker compose -f test-compose.yml up -d
	go test ./...
	docker compose -f test-compose.yml down

