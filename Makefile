.PHONY: run build test tidy smoke clean

run:
	go run .

build:
	go build -o execution-ledger .

test:
	go test ./...

smoke:
	go run ./cmd/smoketest

tidy:
	go mod tidy

clean:
	rm -f execution-ledger *.db *.db-shm *.db-wal
