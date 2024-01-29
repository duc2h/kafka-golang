proto-generate:
	protoc --proto_path=proto proto/kafka/*.proto  --go_out=:pb --go-grpc_out=:pb

up:
	docker compose up -d

down: 
	docker compose down --volumes