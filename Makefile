proto-generate:
	protoc --proto_path=proto proto/kafka/*.proto  --go_out=:pb --go-grpc_out=:pb

up:
	docker compose up

down: 
	docker compose down