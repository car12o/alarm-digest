main:
	@docker-compose up

serve: install.devdeps
	@air -c .air.toml

install.devdeps:
	@GO111MODULE=off go get -v github.com/cosmtrek/air

verify:
	@docker exec -it alarm-digest_nats_1 verify

.PHONY: test
test:
	@go test ./internal/...

e2e.serve:
	@docker-compose up --scale server=3

e2e.test:
	@go test ./test/...
