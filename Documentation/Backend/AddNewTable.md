# Adding a New Backend Table

To add a new backend table, follow these steps:

1. **Define the Model**  
Create the model structure in the `Backend/models/models.go` file.  
Ensure the model includes the necessary fields and annotations for database mapping.

2. **Update Database Migration**  
Add the newly created model to the `Db.AutoMigrate` call in the `Backend/database/database.go` file.  
This ensures the table is automatically created or updated in the database schema.
