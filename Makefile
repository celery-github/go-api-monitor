.PHONY: tidy fmt test run once

tidy:
	go mod tidy

fmt:
	gofmt -w .

test:
	go test ./... -race

run:
	go run ./cmd/monitor --config configs/sample.yaml

once:
	go run ./cmd/monitor --config configs/sample.yaml --once
