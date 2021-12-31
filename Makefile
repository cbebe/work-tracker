PORT := 5000

web: export PORT=$(PORT)
web:
	go run cmd/webserver/main.go

webserver:
	go build cmd/webserver/main.go

cli:
	go build cmd/cli/main.go

USER_ID := $(shell id -u)
GROUP_ID := $(shell id -g)
image:
	docker build \
		--build-arg USER_ID=$(USER_ID) \
		--build-arg GROUP_ID=$(GROUP_ID) \
		--build-arg PORT=$(PORT) \
		-t worktracker .
	
run:
	docker run -p $(PORT):$(PORT) -v $(PWD)/work.db:/app/work.db worktracker

.PHONY: cli webserver web