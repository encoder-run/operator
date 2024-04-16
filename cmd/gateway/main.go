package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/encoder-run/operator/cmd/gateway/middleware"
	"github.com/encoder-run/operator/pkg/graph"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	// Create a router instance
	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{
			// Move these to environment variables.
			"http://localhost:3000",
			"http://localhost:32081",
		},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
		Debug:          true,
	}).Handler)

	// Create a K8sClientManager instance
	km := &middleware.K8sClientManager{}

	// Apply the Kubernetes client middleware
	router.Use(middleware.K8sImpersonationMiddleware(km))

	// GraphQL server setup
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
