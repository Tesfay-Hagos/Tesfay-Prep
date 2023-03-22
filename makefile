postgres:
	sudo docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=tesfay2f -e POSTGRES_PASSWORD=tsionawi@2121 -d postgres:12
createdb:
	sudo docker exec -it postgres12 createdb --username=tesfay2f --owner=tesfay2f simple_bank
dropdb:
	sudo docker exec -it --network bank_netwok postgres12 dropdb --username=tesfay2f simple_bank
migrateup:
	migrate -path simplebank/db/migrations -database 'postgresql://tesfay2f:tsionawi@2121@localhost:5432/simple_bank?sslmode=disable' -verbose up
migrateup1:
	migrate -path simplebank/db/migrations -database 'postgresql://tesfay2f:tsionawi@2121@localhost:5432/simple_bank?sslmode=disable' -verbose up 1
migratedown:
	migrate -path simplebank/db/migrations -database 'postgresql://tesfay2f:tsionawi@2121@localhost:5432/simple_bank?sslmode=disable' -verbose down
migratedown1:
	migrate -path simplebank/db/migrations -database 'postgresql://tesfay2f:tsionawi@2121@localhost:5432/simple_bank?sslmode=disable' -verbose down 1
backtov1:
	migrate -path simplebank/db/migrations -database 'postgresql://tesfay2f:tsionawi@2121@localhost:5432/simple_bank?sslmode=disable' force 1
sqlc:
	sqlc generate
getpid:
	sudo lsof -i :5432
test:
	go test -v ./...
startdocker:
	sudo docker start postgres12
server:
	go run main.go
mock:
	mockgen -package mockdb -destination simplebank/db/mock/store.go tesfayprep/simplebank/db/sqlc Store

.PHONY: createdb dropdb migrateup migratedown getpid sqlc test startdocker server mockgen migrateup1 migratedown1