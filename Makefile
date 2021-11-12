# all our targets are phony (no files to check).
.PHONY: lint clean cron

lint:
	golangci-lint run -c ./golangci.yaml

clean:
	go run main.go clean

cron:
	go run main.go cron