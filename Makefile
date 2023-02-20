## Generate protobuf files
protobuf:
	./run/run-protoc.sh

## Run tests
test:
	./run/test.sh

docker-compose:
	docker-compose up -d gateway-api

end-to-end-tests:
	docker-compose up end-to-end-tests