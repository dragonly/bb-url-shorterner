.PHONY: run test coverage genmock

run:
	go run main.go

genmock:
	mockgen -source=dao/db.go -destination=dao/db_mock.go -package=dao

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
