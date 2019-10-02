package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	graphql_demo "github.com/andreylm/graphql-demo"
	"github.com/andreylm/graphql-demo/api/auth"
	"github.com/andreylm/graphql-demo/api/dal"
	"github.com/andreylm/graphql-demo/api/dataloaders"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := dal.Connect()
	if err != nil {
		panic(err)
	}

	rootHandler := dataloaders.DataloaderMiddleware(
		db,
		handler.GraphQL(
			graphql_demo.NewExecutableSchema(
				graphql_demo.NewRootResolvers(db),
			),
			handler.ComplexityLimit(300),
		),
	)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", auth.Middleware(rootHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
