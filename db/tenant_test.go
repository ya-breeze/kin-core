package db

import (
	"testing"

	"github.com/ya-breeze/kin-core/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	models.TenantModel
	Name string
}

func TestScope(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&TestModel{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Seed data
	db.Create(&TestModel{TenantModel: models.TenantModel{FamilyID: 1}, Name: "Family 1 Data"})
	db.Create(&TestModel{TenantModel: models.TenantModel{FamilyID: 2}, Name: "Family 2 Data"})

	var results []TestModel

	// Test scope for Family 1
	db.Scopes(Scope(1)).Find(&results)
	if len(results) != 1 || results[0].FamilyID != 1 {
		t.Errorf("Scope(1) failed, expected 1 record for Family 1, got %d", len(results))
	}

	// Test scope for Family 2
	db.Scopes(Scope(2)).Find(&results)
	if len(results) != 1 || results[0].FamilyID != 2 {
		t.Errorf("Scope(2) failed, expected 1 record for Family 2, got %d", len(results))
	}
}

func TestCheckOwnership(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&TestModel{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	record := TestModel{TenantModel: models.TenantModel{FamilyID: 1}, Name: "Owned Data"}
	db.Create(&record)

	var found TestModel
	// Correct ownership
	err = CheckOwnership(db, &found, record.ID, 1)
	if err != nil {
		t.Errorf("CheckOwnership failed for correct owner: %v", err)
	}

	// Wrong ownership
	err = CheckOwnership(db, &found, record.ID, 2)
	if err == nil {
		t.Error("CheckOwnership should have failed for wrong owner")
	}
}
