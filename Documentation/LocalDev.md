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