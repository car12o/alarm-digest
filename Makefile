main:
	@docker-compose up

serve: install.devdeps
	@air -c .air.toml

install.devdeps:
	@GO111MODULE=off go get -v github.com/cosmtrek/air
