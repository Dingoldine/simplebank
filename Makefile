clean:
	docker stop postgres 2> /dev/null || true
	docker container prune -f
	fuser -k 5432/tcp
createcontainer:
	docker pull postgres:16-alpine
	docker run --name postgres -d -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:16-alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root simplebank;
dropdb:
	docker exec -it postgres dropdb --username=root  simplebank;
up:
	migrate -verbose -path db/migration/ -database "postgres://root:secret@localhost:5432/simplebank?sslmode=disable" up 
down:
	migrate -path db/migration/ -database "postgres://root:secret@localhost:5432/simplebank?sslmode=disable" down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY:
	all