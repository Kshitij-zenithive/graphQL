package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.62

import (
	"book_management_system/graph/model"
	"book_management_system/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// CreateBook mutation creates a new book record
func (r *mutationResolver) CreateBook(ctx context.Context, input model.CreateBookInput) (*model.Book, error) {
	log.Println("CreateBook called with input:", input)

	// Convert string date to time.Time for database storage
	publishedAt, err := time.Parse("2006-01-02", input.PublishedAt)
	if err != nil {
		log.Println("Error parsing date:", err)
		return nil, errors.New("invalid date format for publishedAt. Expected YYYY-MM-DD")
	}

	// Create new book instance with input data
	book := &models.Book{
		Title:       input.Title,
		Author:      input.Author,
		ISBN:        input.Isbn,
		PublishedAt: publishedAt,
	}

	// Create book record in database
	if err := r.DB.Create(book).Error; err != nil {
		log.Println("Error creating book:", err)
		// Check for ISBN uniqueness violation
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("a book with this ISBN already exists")
		}
		return nil, errors.New("failed to create book: " + err.Error())
	}

	// Convert database model to GraphQL model for response
	return &model.Book{
		ID:          strconv.FormatUint(uint64(book.ID), 10),
		Title:       book.Title,
		Author:      book.Author,
		Isbn:        book.ISBN,
		PublishedAt: book.PublishedAt.Format("2006-01-02"),
		CreatedAt:   book.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   book.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateBook mutation updates an existing book
func (r *mutationResolver) UpdateBook(ctx context.Context, id string, input model.UpdateBookInput) (*model.Book, error) {
	log.Println("UpdateBook called with id:", id, "and input:", input)

	// Find existing book
	var book models.Book
	if err := r.DB.First(&book, id).Error; err != nil {
		log.Println("Error finding book:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, errors.New("database error: " + err.Error())
	}

	// Build updates map for provided fields only
	updates := make(map[string]interface{})
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Author != nil {
		updates["author"] = *input.Author
	}
	if input.Isbn != nil {
		updates["isbn"] = *input.Isbn
	}
	if input.PublishedAt != nil {
		publishedAt, err := time.Parse("2006-01-02", *input.PublishedAt)
		if err != nil {
			log.Println("Error parsing date:", err)
			return nil, errors.New("invalid date format for publishedAt")
		}
		updates["published_at"] = publishedAt
	}

	// Perform update within transaction
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&book).Updates(updates).Error; err != nil {
			log.Println("Error updating book:", err)
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errors.New("a book with this ISBN already exists")
			}
			return err
		}
		// Refresh book data
		return tx.First(&book, id).Error
	})

	if err != nil {
		log.Println("Transaction error:", err)
		return nil, errors.New("failed to update book: " + err.Error())
	}

	// Return updated book as GraphQL model
	return &model.Book{
		ID:          strconv.FormatUint(uint64(book.ID), 10),
		Title:       book.Title,
		Author:      book.Author,
		Isbn:        book.ISBN,
		PublishedAt: book.PublishedAt.Format("2006-01-02"),
		CreatedAt:   book.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   book.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteBook mutation removes a book (soft delete)
func (r *mutationResolver) DeleteBook(ctx context.Context, id string) (bool, error) {
	log.Println("DeleteBook called with id:", id)

	// Attempt to delete the book
	result := r.DB.Delete(&models.Book{}, id)

	if result.Error != nil {
		log.Println("Error deleting book:", result.Error)
		return false, errors.New("failed to delete book: " + result.Error.Error())
	}

	// Check if book was found and deleted
	if result.RowsAffected == 0 {
		log.Println("Book not found")
		return false, errors.New("book not found")
	}

	return true, nil
}

// Books query retrieves all books
func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	log.Println("Books query called")

	var dbBooks []models.Book

	// Retrieve all non-deleted books
	if err := r.DB.Find(&dbBooks).Error; err != nil {
		log.Println("Error fetching books:", err)
		return nil, errors.New("failed to fetch books: " + err.Error())
	}

	// Convert database models to GraphQL models
	fmt.Println(dbBooks)
	books := make([]*model.Book, len(dbBooks))
	for i, book := range dbBooks {
		books[i] = &model.Book{
			ID:          strconv.FormatUint(uint64(book.ID), 10),
			Title:       book.Title,
			Author:      book.Author,
			Isbn:        book.ISBN,
			PublishedAt: book.PublishedAt.Format("2006-01-02"),
			CreatedAt:   book.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   book.UpdatedAt.Format(time.RFC3339),
		}
	}

	return books, nil
}

// Book query retrieves a single book by ID
func (r *queryResolver) Book(ctx context.Context, id string) (*model.Book, error) {
	log.Println("Book query called with id:", id)

	var book models.Book

	// Find book by ID
	if err := r.DB.First(&book, id).Error; err != nil {
		log.Println("Error finding book:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, errors.New("database error: " + err.Error())
	}

	// Convert to GraphQL model
	return &model.Book{
		ID:          strconv.FormatUint(uint64(book.ID), 10),
		Title:       book.Title,
		Author:      book.Author,
		Isbn:        book.ISBN,
		PublishedAt: book.PublishedAt.Format("2006-01-02"),
		CreatedAt:   book.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   book.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }