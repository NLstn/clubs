# PostgreSQL Initialization Files

This directory contains the initialization files for the PostgreSQL database server used in the development container.

## Files

### `init-postgres.sql`
SQL script that runs during PostgreSQL initialization to create all required databases and users for the development environment.

This script creates:

#### Clubs Application Database
- **Database**: `clubs_dev`
- **User**: `clubs_dev`
- **Password**: `clubs_dev_password`

This is the main application database used by the Clubs backend.

#### Keycloak Database
- **Database**: `keycloak_dev`
- **User**: `keycloak_dev`
- **Password**: `keycloak_dev_password`

This is a separate database for Keycloak's authentication data, ensuring the authentication service is isolated from the application data.

## How It Works

When the PostgreSQL Docker container starts for the first time:
1. PostgreSQL looks for `.sql` files in `/docker-entrypoint-initdb.d/`
2. It executes all scripts in alphabetical order
3. The `init-postgres.sql` script creates both databases and their users
4. The initialization only runs once - if the database volume already exists, the script is not re-executed

## Modifying the Configuration

To change database credentials or add new databases:

1. **Edit `init-postgres.sql`** with the desired changes
2. **Remove the existing Docker volume** to trigger re-initialization:
   ```bash
   docker compose down -v
   docker compose up -d
   ```
   ⚠️ **Warning**: This will delete all existing data in both databases!

3. **Rebuild the devcontainer** to apply the changes

Alternatively, you can manually create databases and users by connecting to PostgreSQL:
```bash
# From within the devcontainer
psql -h db -U postgres
```

## Security Note

⚠️ **These credentials are for development only!**

Never use these settings in production:
- Change all passwords to strong, unique values
- Use environment variables or secrets management
- Follow the principle of least privilege for database users
- Enable SSL/TLS for database connections
- Regularly rotate credentials
