# Start all docker services.
start:
	docker compose up -d

# Start only the DB service.
# This is useful for debugging the server.
start-local:
	docker compose up -d db

stop:
	docker compose down --volumes

exec:
	docker exec -it db bash

run-local-client:
	go run client/client.go

# Restart the current DB and then exec into the DB container.
db-restart: stop start wait exec

# Wait sleeps until the DB has been populated.
durations = 1 2 3 4 5 6 7 
wait: 
	$(foreach var,$(durations),sleep 1;)

## Tests
test:
	docker compose run --no-deps --rm server bash -c "cd /go/src && go test -v --race ./..."
	
mock: # requires the installation of mockery on local system: "brew install mockery"
	mockery --all