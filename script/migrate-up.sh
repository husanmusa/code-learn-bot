export POSTGRESQL_URL='postgres://husanmusa:1234@localhost:5432/codelearn?sslmode=disable'

migrate -database ${POSTGRESQL_URL} -path migrations up