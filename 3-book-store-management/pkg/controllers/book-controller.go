package controllers

import (
	"encoding/json"
	"fmt"
	"go-proj/3-book-store-management/pkg/models"
	"go-proj/3-book-store-management/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var NewBook models.Book

func GetBooks(w http.ResponseWriter, r *http.Request) {
	newBooks, err := models.GetAllBooks()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(newBooks)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		http.Error(w, "Book ID required", http.StatusBadRequest)
		return
	}
	// bookId, err := strconv.Atoi(id)
	bookId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("Something went wrong while parsing %v", err), http.StatusInternalServerError)
		return
	}
	book, _ := models.GetBookById(bookId)
	res, err := json.Marshal(book)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {

	newBook := &models.Book{}
	utils.ParseBody(r, newBook)

	book, err := newBook.CreateBook()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(book)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId, ok := params["id"]
	if !ok {
		http.Error(w, "Book ID required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("Something went wrong while parsing %v", err), http.StatusInternalServerError)
		return
	}
	deletedBook, err := models.DeleteBook(id)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(deletedBook)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook models.Book
	utils.ParseBody(r, &updateBook)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Book Id is required", http.StatusBadRequest)
		return
	}
	bookId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("Something went wrong when parsing %v", err), http.StatusInternalServerError)
		return
	}
	// find the book
	var selectedBook models.Book
	selectedBooks, db := models.GetBookById(bookId)
	if len(selectedBooks) != 0 {
		selectedBook = selectedBooks[0]
		if updateBook.Name != "" {
			selectedBook.Name = updateBook.Name
		}
		if updateBook.Author != "" {
			selectedBook.Author = updateBook.Author
		}
		if updateBook.Publication != "" {
			selectedBook.Publication = updateBook.Publication
		}
		db.Save(selectedBook)
	}
	res, _ := json.Marshal(selectedBook)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookByName(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	bookName, ok := vars["name"]
	if !ok {
		http.Error(w, "Book name must be included", http.StatusBadRequest)
		return
	}
	book, err := models.GetBookByName(bookName)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
