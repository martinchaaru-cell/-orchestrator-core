.PHONY: help build run test clean docker-build

help:
	@echo "Orchestrator Core - Modular Security Task Distribution"
	@echo ""
	@echo "Commands:"
	@echo "  make build        - Build Go binaries"
	@echo "  make run          - Run with docker-compose"
	@echo "  make test         - Run tests"
	@echo "  make docker-build - Build Docker image"
	@echo "  make clean        - Clean build artifacts"

build:
	@echo "Building orchestrator..."
	cd orchestrator && go build -o ../bin/orchestrator .
	@echo "Building agent..."
	cd agent && go build -o ../bin/agent .
	@echo "✅ Binaries ready in bin/"

run:
	docker-compose up

test:
	cd orchestrator && go test ./... || true
	cd agent && go test ./... || true

docker-build:
	docker build -f orchestrator/Dockerfile -t orchestrator:latest orchestrator/
	docker build -f agent/Dockerfile -t agent:latest agent/

clean:
	rm -rf bin/
	docker-compose down
