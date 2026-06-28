package models

import "testing"

// TestAcquisitionTargetUnique verifies the catalog spine migrates and that a
// Book may have at most one AcquisitionTarget per BookType.
func TestAcquisitionTargetUnique(t *testing.T) {
	db, err := ConnectTestDB()
	if err != nil {
		t.Fatalf("ConnectTestDB: %v", err)
	}

	author := Author{Name: "Brandon Sanderson"}
	if err := db.Create(&author).Error; err != nil {
		t.Fatalf("create author: %v", err)
	}

	book := Book{Authors: []Author{author}, Title: "Oathbringer"}
	if err := db.Create(&book).Error; err != nil {
		t.Fatalf("create book: %v", err)
	}

	var ebook BookType
	if err := db.Where("name = ?", "ebook").First(&ebook).Error; err != nil {
		t.Fatalf("find ebook type: %v", err)
	}

	if err := db.Create(&AcquisitionTarget{BookID: book.ID, BookTypeID: ebook.ID}).Error; err != nil {
		t.Fatalf("create first target: %v", err)
	}
	if err := db.Create(&AcquisitionTarget{BookID: book.ID, BookTypeID: ebook.ID}).Error; err == nil {
		t.Fatal("expected unique-constraint error on duplicate (book, booktype), got nil")
	}
}
