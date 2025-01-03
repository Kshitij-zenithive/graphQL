package main

import (
	"book_management_system/graph"
	"book_management_system/internal/database"
	"log"
	"net/http"
	

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	// Load environment variables from .env file
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file:", err)
	// }

	// Initialize database connection
	if err := database.InitDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Get port from environment or use default
	port := "8080"
	if port == "" {
		port = defaultPort
	}

	// Create resolver with database connection
	resolver := graph.NewResolver(database.DB)

	// Create GraphQL server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	// Configure server
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Set up HTTP handlers
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Start server
	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
