

docker-build:
	docker build -t go-authenticator .

docker-run: docker-build
	docker run --rm --name go-authenticator -p 9123:9123 --network=go-authenticator_auth go-authenticator

run:
	go run main.go
