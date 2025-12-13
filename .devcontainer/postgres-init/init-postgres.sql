-- PostgreSQL Initialization Script for Clubs Development Environment
-- This script runs automatically when the PostgreSQL container starts for the first time.
-- It creates both the Clubs application database and the Keycloak database.

-- Create Clubs application database and user
CREATE USER clubs_dev WITH PASSWORD 'clubs_dev_password';
CREATE DATABASE clubs_dev OWNER clubs_dev;
GRANT ALL PRIVILEGES ON DATABASE clubs_dev TO clubs_dev;

-- Create Keycloak database and user
CREATE USER keycloak_dev WITH PASSWORD 'keycloak_dev_password';
CREATE DATABASE keycloak_dev OWNER keycloak_dev;
GRANT ALL PRIVILEGES ON DATABASE keycloak_dev TO keycloak_dev;
