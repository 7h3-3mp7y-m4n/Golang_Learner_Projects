package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/7h3-3mp7y-m4n/booker/store"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "postgres://sslmode=disable" //replace it with db postgres config
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", index(db)).Methods("GET")
	r.HandleFunc("/create", createBook(db)).Methods("GET", "POST")
	r.HandleFunc("/update", updateHandler(db)).Methods("GET", "POST")
	r.HandleFunc("/delete", deleteHandler(db)).Methods("POST")

	fmt.Println("Starting the server on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func index(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := store.New(db)
		ctx := context.Background()
		books, err := queries.ListBooks(ctx)
		if err != nil {
			http.Error(w, "Failed to list books", http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, books)
	}
}

func createBook(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			queries := store.New(db)
			ctx := context.Background()
			price := 100
			err := queries.CreateBook(ctx, store.CreateBookParams{
				Title:  r.FormValue("title"),
				Author: r.FormValue("author"),
				Genre:  r.FormValue("genre"),
				Price:  price,
			})
			if err != nil {
				http.Error(w, "Failed to create book", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		tmpl := template.Must(template.ParseFiles("book_form.html"))
		tmpl.Execute(w, nil)
	}
}

func updateHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := store.New(db)
		ctx := context.Background()

		if r.Method == http.MethodPost {
			id, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				http.Error(w, "Invalid book ID", http.StatusBadRequest)
				return
			}

			price, err := strconv.Atoi(r.FormValue("price"))
			if err != nil {
				http.Error(w, "Invalid price", http.StatusBadRequest)
				return
			}

			err = queries.UpdateBook(ctx, store.UpdateBookParams{
				ID:     int32(id),
				Title:  r.FormValue("title"),
				Author: r.FormValue("author"),
				Genre:  r.FormValue("genre"),
				Price:  price,
			})
			if err != nil {
				http.Error(w, "Failed to update book", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		book, err := queries.GetBookByID(ctx, int32(id))
		if err != nil {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}

		tmpl := template.Must(template.ParseFiles("book_update.html"))
		tmpl.Execute(w, book)
	}
}

func deleteHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			queries := store.New(db)
			ctx := context.Background()

			idStr := r.FormValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid book ID", http.StatusBadRequest)
				return
			}

			err = queries.DeleteBook(ctx, int32(id))
			if err != nil {
				http.Error(w, "Failed to delete book", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}
