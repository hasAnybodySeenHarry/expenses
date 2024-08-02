fmt:
	@echo "Formatting code..."
	go fmt ./...

run: fmt
	@echo "Running ./cmd/api/..."
	go run ./cmd/api/

.PHONY: fmt run