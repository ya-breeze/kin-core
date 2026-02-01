# kin-core

A shared library for multi-tenant (family-based) Go applications.

## Features

- **Standardized Identity**: Shared `Family` and `User` models.
- **Data Isolation**: `TenantModel` base and GORM scopes to enforce cross-family privacy.
- **Authentication**: JWT generation and Gin middleware for automated context injection.

## Installation

```bash
go get github.com/ya-breeze/kin-core
```

## Usage

### 1. Define Your Models

Embed `models.TenantModel` to make any struct tenant-aware.

```go
import "github.com/ya-breeze/kin-core/models"

type ShoppingList struct {
    models.TenantModel
    Title string `json:"title"`
}
```

### 2. Configure Gin Middleware

Register the middleware to automatically handle JWTs and inject `family_id` into requests.

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/ya-breeze/kin-core/middleware"
)

func main() {
    r := gin.Default()
    secret := []byte("your-jwt-secret")
    
    // Injects 'user_id' and 'family_id' into context
    r.Use(middleware.AuthMiddleware(secret))
    
    // ... routes
}
```

### 3. Enforce Isolation in Queries

Use `db.Scope` in your handlers to ensure users only see data belonging to their family.

```go
import "github.com/ya-breeze/kin-core/db"

func GetLists(c *gin.Context) {
    familyID := c.MustGet("family_id").(uint)
    var lists []models.ShoppingList
    
    // Scopes(db.Scope(familyID)) adds WHERE family_id = ?
    if err := database.Scopes(db.Scope(familyID)).Find(&lists).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch"})
        return
    }
    
    c.JSON(200, lists)
}
```

### 4. Verify Ownership for Mutations

Before updating or deleting, verify the record belongs to the user's family.

```go
func DeleteList(c *gin.Context) {
    familyID := c.MustGet("family_id").(uint)
    id := c.Param("id")
    
    var list models.ShoppingList
    if err := db.CheckOwnership(database, &list, id, familyID); err != nil {
        c.JSON(404, gin.H{"error": "Not found or unauthorized"})
        return
    }
    
    database.Delete(&list)
}
```
