# Local Development

For local development you need to replicate the 3 components of clubs in your local environment.

## Database

For the database, you have to do the following steps:

1. Create postgresql docker container
```bash
docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres:latest
```

2. Exec into the created container and run psql as user postgres
```bash
docker exec -it -u postgres postgres psql
```

3. Create a new user named clubs
```sql
CREATE USER clubs WITH PASSWORD 'yourpassword';
```

4. Create a new database also named clubs and set the user clubs as its owner, so it can create tables and stuff while automigrating the schema
```sql
CREATE DATABASE clubs OWNER clubs;
```

## Backend

The backend can be run standalone as well if the database is up at the place defined in .env. This file is only relevant in development.

Simply run the backend like this.
```bash
go run main.go
```

## Frontend

The frontend will also work out of the box if the backend is up and running, but you can choose which backend to use. Change the VITE_API_HOST in .env.development to the host of the backend and run

```bash
npm run dev
```

## Development Authentication

For testing and development purposes (including AI agent workflows), the application provides a development-only authentication endpoint that bypasses email verification.

### Setup

1. Ensure `ENABLE_DEV_AUTH=true` is set in your Backend/.env file
2. Set a `JWT_SECRET` in your Backend/.env file (any string will work for development)
3. Start the backend server

### Usage

Authenticate with any email address without verification:

```bash
# Login and get tokens
curl -X POST http://localhost:8080/api/v1/auth/dev-login \
  -H "Content-Type: application/json" \
  -d '{"email": "dev@example.com"}'

# Response will include access and refresh tokens:
# {
#   "access": "eyJhbGc...",
#   "refresh": "eyJhbGc...",
#   "profileComplete": false
# }

# Use the access token for authenticated requests:
curl -X GET http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**⚠️ Security Note:** This endpoint is for development only and will return 404 in production environments where `ENABLE_DEV_AUTH` is not set to `true`.