#!/bin/bash
set -e

# Create databases for each service
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    CREATE DATABASE auth_db;
    CREATE DATABASE task_db;
    CREATE DATABASE project_db;
    CREATE DATABASE notification_db;
EOSQL

echo "Databases created successfully!"
