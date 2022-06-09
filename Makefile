PORT := 5000

web:
	PORT=$(PORT) go run cmd/web_server/main.go
runbot:
	go run cmd/bot/main.go

web_server: bin/web_server
cli: bin/cli
bot: bin/bot
bin/%: cmd/%/main.go
	go build $<
	mv main $@

USER_ID := $(shell id -u)
GROUP_ID := $(shell id -g)
image:
	docker build \
		--build-arg USER_ID=$(USER_ID) \
		--build-arg GROUP_ID=$(GROUP_ID) \
		--build-arg PORT=$(PORT) \
		-t worktracker .
	
watch: bin/bot
	find cmd/ pkg/ .env | entr sh -c '. ./.env && make runbot'

	
run:
	docker run -p $(PORT):$(PORT) -v $(PWD)/work.db:/app/work.db worktracker

.PHONY: cli webserver web image run
