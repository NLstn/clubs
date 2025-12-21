<div align="center">
  <img src="../assets/logo.png" alt="Clubs Logo" width="150"/>
  
  # Adding a New Backend Table
  
  **Guide for extending the database schema**
</div>

---

## 📋 Overview

This guide walks you through the process of adding a new database table to the Clubs backend application.

## 🔧 Steps to Add a New Table

### 1. Define the Model

Create the model structure in the `Backend/models/` directory (typically in a new file or in an existing related model file).

**Important**: Ensure the model includes:
- Proper GORM tags for database mapping
- JSON tags in **PascalCase** for OData v2 API compatibility
- Appropriate field types and constraints

**Example**:
```go
// Backend/models/your_model.go
type YourModel struct {
    ID        string    `json:"ID" gorm:"type:uuid;primaryKey"`
    Name      string    `json:"Name" gorm:"not null"`
    CreatedAt time.Time `json:"CreatedAt" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"UpdatedAt" gorm:"autoUpdateTime"`
}
```

### 2. Update Database Migration

Add the newly created model to the `AutoMigrate` call in the `Backend/database/database.go` file.

**Example**:
```go
// Backend/database/database.go
func InitializeDatabase() error {
    // ... existing code ...
    
    err := Db.AutoMigrate(
        &models.Club{},
        &models.Member{},
        &models.YourModel{}, // Add your new model here
        // ... other models ...
    )
    
    // ... rest of the code ...
}
```

This ensures the table is automatically created or updated in the database schema when the application starts.

## ⚠️ Important Considerations

### OData v2 Compatibility
**All JSON field names must use PascalCase** for OData v2 API endpoints:
```go
// ✅ Correct
type User struct {
    ID        string `json:"ID"`
    FirstName string `json:"FirstName"`
}

// ❌ Wrong
type User struct {
    ID        string `json:"id"`
    FirstName string `json:"firstName"`
}
```

### Foreign Keys
When adding relationships, use proper foreign key constraints:
```go
type YourModel struct {
    ClubID string `json:"ClubID" gorm:"type:uuid;not null"`
    Club   Club   `json:"Club" gorm:"foreignKey:ClubID"`
}
```

### Indexes
Add indexes for frequently queried fields:
```go
type YourModel struct {
    Email string `json:"Email" gorm:"uniqueIndex"`
    Name  string `json:"Name" gorm:"index"`
}
```

## 🧪 Testing

After adding a new table:

1. **Build the backend**:
   ```bash
   cd Backend
   go build
   ```

2. **Run the application** to trigger auto-migration:
   ```bash
   go run main.go
   ```

3. **Verify the table** was created:
   ```sql
   -- Connect to your database
   \dt  -- List all tables
   \d your_models  -- Describe your new table
   ```

## 📚 Related Documentation

- [API Documentation](API.md)
- [GORM Documentation](https://gorm.io/docs/)
- [OData v2 Guidelines](../Frontend/README.md)

---

**Note**: Remember to add appropriate API endpoints and handlers after creating the database table. See [API.md](API.md) for endpoint documentation patterns.
