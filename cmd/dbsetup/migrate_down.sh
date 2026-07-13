#!/bin/bash
# WIP: use this to migrate down, since migrate up happens at docker exec 

# Load .env
set -a

if [ -f ".env" ]; then
    . "./.env"

else
    echo "Error: .env file not found in the parent directory." >&2
    exit 1
fi

set +a

POSTGRES_URL=postgres://$DB_USERNAME:$DB_PASSWORD@localhost:$DB_HOST_PORT/$DB_NAME?sslmode=disable
migrate -database ${POSTGRES_URL} -path cmd/dbsetup/migration_files down

if [ $? -eq 0 ]; then
    echo "Successfully reverse migrated."
else
    echo "Something went wrong."
fi
