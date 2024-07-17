test:
	docker compose -f test-compose.yml up -d ; sleep 1
	go test ./... -v ;\
	docker compose -f test-compose.yml down

