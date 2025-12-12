-- Create Keycloak database and user
CREATE USER keycloak_dev WITH PASSWORD 'keycloak_dev_password';
CREATE DATABASE keycloak_dev OWNER keycloak_dev;
GRANT ALL PRIVILEGES ON DATABASE keycloak_dev TO keycloak_dev;
