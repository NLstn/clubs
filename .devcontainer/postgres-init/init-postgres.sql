-- PostgreSQL Initialization Script for Civo Development Environment
-- This script runs automatically when the PostgreSQL container starts for the first time.
-- It creates both the Civo application database and the Keycloak database.

-- Create Civo application database and user
CREATE USER civo_dev WITH PASSWORD 'civo_dev_password';
CREATE DATABASE civo_dev OWNER civo_dev;
GRANT ALL PRIVILEGES ON DATABASE civo_dev TO civo_dev;

-- Create Keycloak database and user
CREATE USER keycloak_dev WITH PASSWORD 'keycloak_dev_password';
CREATE DATABASE keycloak_dev OWNER keycloak_dev;
GRANT ALL PRIVILEGES ON DATABASE keycloak_dev TO keycloak_dev;
