# Create Migration Commands
c_m:
	@echo "Creating migrations..."
	migrate create -ext sql -dir db/migrations -seq $(name)

# PostgreSQL Commands
p_up:
	@echo "Starting PostgreSQL..."
	docker-compose up -d

p_down:
	@echo "Stopping and removing PostgreSQL..."
	docker-compose down

db_up:
	@echo "Creating a database..."
	docker exec -it pooler_finTech_postgres createdb --username=root --owner=root pooler_fintech_postgres_db


db_down:
	@echo "Dropping the database..."
	docker exec -it pooler_finTech_postgres dropdb --username=root pooler_fintech_postgres_db


# Apply Migration Up and Down Commands
m_up:
	@echo "Applying database migrations (up)..."
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/pooler_fintech_postgres_db?sslmode=disable" up


m_down:
	@echo "Reverting database migrations (down)..."
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/pooler_fintech_postgres_db?sslmode=disable" down


# SQLC Commands
sqlc:
	@echo "Generating SQLC code..."
	sqlc generate

# Testing
test:
	@echo "Running all tests with coverage..."
	go test -v -cover ./...

# package for JWT token
jwt:
	go get github.com/golang-jwt/jwt
# Start Development Server
start:
	@echo "Starting the development server..."
	CompileDaemon -command = "./pooler_Remmitance_Application"

#creat_db
#CREATE DATABASE pooler_fintech_postgres_db;
#CREATE USER root WITH PASSWORD 'secret';
#GRANT ALL PRIVILEGES ON DATABASE pooler_fintech_postgres_db TO root;
#GRANT ALL PRIVILEGES ON SCHEMA public TO root;

