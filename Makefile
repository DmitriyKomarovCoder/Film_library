include .env
export

db:
	docker exec -it hammy-db psql -U $(DB_USER) -d $(DB_NAME)

doc:
	swag init -g cmd/app/main.go

cover:
	sh scripts/coverage_test.sh

