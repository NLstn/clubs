#!/bin/bash
set -e

# Create the clubs user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER clubs WITH PASSWORD 'yourpassword';
    CREATE DATABASE clubs OWNER clubs;
    GRANT ALL PRIVILEGES ON DATABASE clubs TO clubs;
EOSQL

echo "Database 'clubs' and user 'clubs' created successfully"