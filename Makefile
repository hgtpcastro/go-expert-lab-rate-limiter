.PHONY: start
start:
	@docker compose -f deployments/docker-compose/docker-compose.yaml up --build -d	

.PHONY: stop
stop:
	@docker compose -f deployments/docker-compose/docker-compose.yaml stop

.PHONY: clean
clean: stop
	@docker compose -f deployments/docker-compose/docker-compose.yaml down -v --remove-orphans	

.PHONY: test
test:	
	@go test -v ./...

.PHONY: test-smoke
test-smoke:	
	@docker compose -f deployments/docker-compose/docker-compose.yaml run k6 run /scripts/smoke-test.js	

.PHONY: test-stress
test-stress:	
	@docker compose -f deployments/docker-compose/docker-compose.yaml run k6 run /scripts/stress-test.js	