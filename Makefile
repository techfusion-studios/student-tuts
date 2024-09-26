proto:
	protoc --proto_path=pb --go_out=./pb --go-grpc_out=./pb --grpc-gateway=pb protos/*.proto

clean:
	rm -rf pb docs

buf:
	buf generate

start:
	go run cmd/main.go

watch:
	air

test:
	go test ./...

docker-build:
	docker build -t techfusion/student-api:1.0.0 .

docker-run:
	docker run -e DB_NAME=students_db -e DB=sqlite -e APP.PORT=9000 -e APP.GRPC_PORT=50051 techfusion/student-api:1.0.0

docker:
	make docker-build
	make docker-run