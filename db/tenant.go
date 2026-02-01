package db

import "gorm.io/gorm"

// Scope returns a GORM scope that filters results by FamilyID
func Scope(familyID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("family_id = ?", familyID)
	}
}

// CheckOwnership verifies that a record with the given ID exists and belongs to the specified FamilyID.
// It populates the model interface with the found record.
func CheckOwnership(db *gorm.DB, model interface{}, id interface{}, familyID uint) error {
	return db.Where("id = ? AND family_id = ?", id, familyID).First(model).Error
}
