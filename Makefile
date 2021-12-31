web: export PORT=8080
web:
	go run cmd/webserver/main.go

webserver:
	go build cmd/webserver/main.go

cli:
	go build cmd/cli/main.go

.PHONY: cli webserver web