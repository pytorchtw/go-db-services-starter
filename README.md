# go-db-services-starter

# start up docker postgresql db server
# this will create a test database directory at ./db_data/test_db if not exists
docker-compose up postgresql

# create db migration scripts
migrate create -ext sql -dir db_data/migrations -seq create_pages_table

export POSTGRESQL_URL='postgres://postgres:postgres123@localhost:5432/docker?sslmode=disable'
migrate -database ${POSTGRESQL_URL} -path db_data/migrations up
migrate -database ${POSTGRESQL_URL} -path db_data/migrations down
migrate -database ${POSTGRESQL_URL} -path db_data/migrations force VERSION

# generate the sqlboiler db models code
sqlboiler psql

# test the generated model
go test models/*.go -v

# test all
go test ./... -v
