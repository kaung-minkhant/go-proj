package models

import (
	"go-proj/3-book-store-management/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `gorm:"" json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() (*Book, error) {
	result := db.Create(&b)
	if result.Error != nil {
		return nil, result.Error
	}
	return b, nil
}

func GetAllBooks() ([]Book, error) {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

func GetBookById(Id int64) ([]Book, *gorm.DB) {
	var getBook []Book
	db := db.Where("ID = ?", Id).Find(&getBook)
	return getBook, db
}

func DeleteBook(Id int64) (*Book, error) {
	var book Book
	result := db.Where("ID = ?", Id).Delete(&book)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func GetBookByName(name string) ([]Book, error) {
	var book []Book
	result := db.Where("name = ?", name).Find(&book)
	if result.Error != nil {
		return nil, result.Error
	}
	return book, nil
}
